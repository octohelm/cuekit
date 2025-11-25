package godebug

import (
	"sync"
)

import (
	_ "unsafe"
)

// Note: Be careful about new imports here. Any package
// that internal/godebug imports cannot itself import internal/godebug,
// meaning it cannot introduce a GODEBUG setting of its own.
// We keep imports to the absolute bare minimum.

// A Setting is a single setting in the $GODEBUG environment variable.
type Setting struct {
	name string
	once sync.Once
}

func New(name string) *Setting {
	return &Setting{name: name}
}

// Name returns the name of the setting.
func (s *Setting) Name() string {
	return s.name
}

func (s *Setting) Undocumented() bool {
	return s.name != "" && s.name[0] == '#'
}

func (s *Setting) String() string {
	return s.Name() + "=" + s.Value()
}

func (s *Setting) IncNonDefault() {
}

func (s *Setting) Value() string {
	return ""
}
