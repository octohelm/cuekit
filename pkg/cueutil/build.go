package cueutil

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/parser"
)

func BuildFile(src any) (cue.Value, error) {
	b, err := parser.ParseFile("x.cue", src)
	if err != nil {
		return cue.Value{}, err
	}
	return cuecontext.New().BuildFile(b), nil
}
