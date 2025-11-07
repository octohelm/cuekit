package cuecontext

import (
	"errors"
	"path/filepath"

	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/load"

	"github.com/octohelm/cuekit/pkg/mod/module"
)

type Instance = build.Instance

func BuildInstance(c *load.Config, inputs []string) (*Instance, error) {
	if len(inputs) == 0 {
		files, err := module.WalkCueFile(c.Dir, ".")
		if err != nil {
			return nil, err
		}
		inputs = files
	}

	files := make([]string, len(inputs))

	for i, f := range inputs {
		if filepath.IsAbs(f) {
			rel, _ := filepath.Rel(c.Dir, f)
			files[i] = "./" + rel
		} else {
			files[i] = f
		}
	}

	insts := load.Instances(files, c)
	if len(insts) != 1 {
		return nil, errors.New("invalid instance")
	}

	if err := insts[0].Err; err != nil {
		return nil, err
	}
	return insts[0], nil
}
