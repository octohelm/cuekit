package cuecontext

import (
	"iter"
	"slices"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

func Build(c *Config, input iter.Seq[string], opts ...cuecontext.Option) (cue.Value, error) {
	inst, err := BuildInstance(c, slices.Collect(input))
	if err != nil {
		return cue.Value{}, err
	}

	v := cuecontext.New(opts...).BuildInstance(inst)
	if err := v.Err(); err != nil {
		return cue.Value{}, err
	}
	return v, nil
}

func EvalJSON(c *Config, inputs ...string) ([]byte, error) {
	v, err := Build(c, slices.Values(inputs))
	if err != nil {
		return nil, err
	}
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return v.MarshalJSON()
}
