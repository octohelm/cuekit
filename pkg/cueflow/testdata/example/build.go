package example

import (
	"context"
	"fmt"

	"cuelang.org/go/cue"
	"github.com/octohelm/cuekit/pkg/cueflow"
	"github.com/octohelm/cuekit/pkg/cueflow/task"
	"github.com/octohelm/cuekit/pkg/cuepath"
)

func init() {
	Registry.Register(&Build{})
}

type Build struct {
	task.Group

	Input  WorkDir                `json:"input"`
	Steps  []BuildDoStepInterface `json:"steps"`
	Output WorkDir                `json:"-" output:"output"`
}

type BuildDoStepInterface struct {
	Input  WorkDir `json:"input,omitzero"`
	Output WorkDir `json:"output"`
}

func (x *Build) Do(ctx context.Context) error {
	t := x.T()

	path := cuepath.Join(t.Path(), cue.ParsePath("input"))

	if err := t.Scope().DecodePath(path, &x.Input); err != nil {
		return err
	}

	step := &BuildDoStepInterface{}
	step.Output = x.Input

	for stepValue, err := range task.Steps(t.Scope().LookupPath(t.Path())) {
		if err != nil {
			return err
		}

		stepPath := stepValue.Path()

		if err := t.Scope().FillPath(stepPath, map[string]any{"input": step.Output}); err != nil {
			return err
		}

		if err := cueflow.RunSubTasks(ctx, t.Scope(), func(p cue.Path) (bool, cue.Path) {
			if cuepath.Prefix(p, stepPath) {
				return true, cuepath.TrimPrefix(p, stepPath)
			}
			return false, cue.MakePath()
		}); err != nil {
			return err
		}

		step = &BuildDoStepInterface{}
		if err := t.Scope().DecodePath(stepPath, step); err != nil {
			return fmt.Errorf("decode result failed: %w", err)
		}
	}

	x.Output = step.Output

	return nil
}
