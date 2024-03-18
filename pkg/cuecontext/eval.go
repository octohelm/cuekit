package cuecontext

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

func Build(c *Config, inputs ...string) (cue.Value, error) {
	inst, err := BuildInstance(c, inputs)
	if err != nil {
		return cue.Value{}, err
	}

	v := cuecontext.New().BuildInstance(inst)
	if err := v.Err(); err != nil {
		return cue.Value{}, err
	}
	return v, nil
}

func EvalJSON(c *Config, inputs ...string) ([]byte, error) {
	v, err := Build(c, inputs...)
	if err != nil {
		return nil, err
	}
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return v.MarshalJSON()
}
