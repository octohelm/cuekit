package task

import (
	"context"

	"github.com/octohelm/cuekit/pkg/cueflow"
)

type Checkpoint struct {
	// no need the ok
	Task `json:"-"`
}

var _ cueflow.Checkpoint = &Checkpoint{}

func (Checkpoint) AsCheckpoint() bool {
	return true
}

func (Checkpoint) Do(ctx context.Context) error {
	return nil
}
