package cuepath

import (
	"cuelang.org/go/cue"
	"github.com/go-json-experiment/json/jsontext"
)

func AsJSONPointer(p cue.Path) jsontext.Pointer {
	var pt jsontext.Pointer
	for _, s := range p.Selectors() {
		if s.Type() == cue.StringLabel {
			pt = pt.AppendToken(s.Unquoted())
			continue
		}
		pt = pt.AppendToken(s.String())
	}
	return pt
}
