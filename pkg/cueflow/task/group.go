package task

import "github.com/octohelm/cuekit/pkg/cueflow"

type Group struct {
	Task

	t cueflow.Task
}

var _ cueflow.TaskUnmarshaler = &Group{}

func (v *Group) UnmarshalTask(t cueflow.Task) error {
	v.t = t
	return nil
}

func (v *Group) T() cueflow.Task {
	return v.t
}
