package example

import (
	"context"
	"fmt"
	"github.com/octohelm/cuekit/pkg/cueflow/task"
)

func init() {
	Registry.Register(&Exec{})
}

type Exec struct {
	task.Task

	Cwd WorkDir `json:"cwd"`

	Command []string          `json:"cmd"`
	Env     map[string]string `json:"env,omitzero"`

	Stdout string `json:"-" output:"stdout,omitzero"`
}

func (local *Exec) Do(ctx context.Context) error {
	fmt.Println("EXEC", local.Command)

	local.Stdout = fmt.Sprintf("%s %s %v", local.Cwd.Ref.ID, local.Command, local.Env)
	return nil
}
