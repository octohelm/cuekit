package example

import (
	"context"
	
	"github.com/octohelm/cuekit/pkg/cueflow/task"
)

func init() {
	Registry.Register(&Local{})
}

type Local struct {
	task.Task

	Source  string  `json:"source" default:"."`
	WorkDir WorkDir `json:"-" output:"dir"`
}

func (local *Local) Do(ctx context.Context) error {
	local.WorkDir.Ref.ID = "local:"
	return nil
}
