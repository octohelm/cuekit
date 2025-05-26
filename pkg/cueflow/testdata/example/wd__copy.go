package example

import (
	"context"
	"github.com/octohelm/cuekit/pkg/cueflow/task"
)

func init() {
	Registry.Register(&Copy{})
}

type Copy struct {
	task.Task

	Input WorkDir `json:"input"`

	Output WorkDir `json:"-" output:"output"`
}

func (c *Copy) Do(ctx context.Context) error {
	c.Output = c.Input
	return nil
}
