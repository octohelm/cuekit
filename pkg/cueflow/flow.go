package cueflow

var _ FlowTask = TaskImpl{}

type TaskImpl struct{}

func (TaskImpl) flowTask() {
}

var _ FlowControl = FlowControlImpl{}

type FlowControlImpl struct{}

func (FlowControlImpl) flowTask() {
}

func (FlowControlImpl) flowControl() {
}

type FlowTask interface {
	flowTask()
}

type FlowControl interface {
	FlowTask

	flowControl()
}

type BeforeAll interface {
	BeforeAll() bool
}

func IsBeforeAll(t Task) bool {
	return t.Value().LookupPath(TaskBeforeAllPath).Exists()
}

type CanSkip interface {
	Skip() bool
}

type CacheDisabler interface {
	CacheDisabled() bool
}

type Checkpoint interface {
	AsCheckpoint() bool
}

type Successor interface {
	Success() bool
}

type ResultValuer interface {
	ResultValue() map[string]any
}

type TaskFeedback interface {
	Done(err error)
}
