package task

import (
	"github.com/octohelm/x/ptr"

	"github.com/octohelm/cuekit/pkg/cueflow"
)

type Task struct {
	// task hook to make task could run after some others
	Dep any `json:"$dep,omitzero"`
	// task result
	Ok *bool `json:"-" output:"$ok,omitzero"`

	cueflow.TaskImpl
}

var _ cueflow.Successor = &Task{}

func (t *Task) Success() bool {
	return t.Ok != nil && *t.Ok
}

var _ cueflow.TaskFeedback = &Task{}

func (t *Task) Done(err error) {
	if t.Ok == nil {
		t.Ok = ptr.Ptr(err == nil)
	}
}

type SetupTask struct {
	Task
}

var _ cueflow.BeforeAll = &SetupTask{}

func (SetupTask) BeforeAll() bool {
	return true
}
