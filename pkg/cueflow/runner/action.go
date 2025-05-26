package runner

import (
	"context"

	"github.com/octohelm/cuekit/pkg/cueflow"
)

func AsAction(r Registry) cueflow.Action {
	return func(ctx context.Context, task cueflow.Task) error {
		runner, err := r.New(task)
		if err != nil {
			return err
		}
		return runner.Run(ctx)
	}
}
