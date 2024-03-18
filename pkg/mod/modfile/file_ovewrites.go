package modfile

import (
	"fmt"
	"strings"
	"sync"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/mod/module"
)

type FileOverwrites struct {
	Deps map[string]*DepOverwrite `json:"deps,omitempty"`

	// Path resolved
	Path string `json:"path,omitempty"`
	// Version resolved
	Version string `json:"version,omitempty"`

	mu sync.RWMutex
}

func (f *FileOverwrites) AddDep(mpath string, dep *DepOverwrite) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.Deps == nil {
		f.Deps = map[string]*DepOverwrite{}
	}
	f.Deps[mpath] = dep
}

func (f *FileOverwrites) GetDep(mpath string) (*DepOverwrite, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.Deps) == 0 {
		return nil, false
	}

	dep, ok := f.Deps[mpath]
	if ok {
		return dep, true
	}

	for m, dep := range f.Deps {
		base, _, ok := module.SplitPathVersion(m)
		if ok && base == mpath {
			return dep, true
		}
	}

	return nil, false
}

func (f *FileOverwrites) Format() ([]byte, error) {
	v := cuecontext.New().Encode(f)
	if err := v.Validate(cue.Concrete(true)); err != nil {
		return nil, err
	}
	n := v.Syntax(cue.Concrete(true)).(*ast.StructLit)
	data, err := format.Node(&ast.File{
		Decls: n.Elts,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot format: %v", err)
	}
	return data, err
}

func (f *FileOverwrites) Load(data []byte) error {
	v := cuecontext.New().CompileBytes(data)
	if err := v.Validate(); err != nil {
		return err
	}
	if err := v.Decode(f); err != nil {
		return err
	}
	return nil
}

func (f *FileOverwrites) IsZero() bool {
	return f == nil || (f.Version == "") && len(f.Deps) == 0
}

type DepOverwrite struct {
	Path    string `json:"path,omitempty"`
	Source  string `json:"source,omitempty"`
	Version string `json:"v,omitempty"`
}

func (o *DepOverwrite) IsLocalReplacement() bool {
	return o.Path != "" && strings.HasPrefix(o.Path, ".")
}
