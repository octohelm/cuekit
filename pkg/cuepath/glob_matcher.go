package cuepath

import (
	"fmt"

	"cuelang.org/go/cue"

	"github.com/gobwas/glob"
)

type GlobMatcher interface {
	Match(p cue.Path) bool
}

func CompileGlobMatcher(rules ...string) (GlobMatcher, error) {
	m := &matcher{}

	for _, x := range rules {
		p := cue.ParsePath(x)
		if err := p.Err(); err != nil {
			return nil, fmt.Errorf("compile %q failed: %w", x, err)
		}

		r, err := glob.Compile(string(AsJSONPointer(p)))
		if err != nil {
			return nil, fmt.Errorf("compile %q failed: %w", x, err)
		}

		m.rules = append(m.rules, r)
	}

	return m, nil
}

type matcher struct {
	rules []glob.Glob
}

func (m *matcher) Match(p cue.Path) bool {
	if m.rules == nil {
		return false
	}

	targetPath := string(AsJSONPointer(p))
	for _, r := range m.rules {
		if r.Match(targetPath) {
			return true
		}
	}
	return false
}
