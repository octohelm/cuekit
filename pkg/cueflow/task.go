package cueflow

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"os"
	"slices"
	"sync"
	"sync/atomic"

	"cuelang.org/go/cue"
	"github.com/octohelm/cuekit/pkg/cueflow/graph"
	"github.com/octohelm/cuekit/pkg/cuepath"
	contextx "github.com/octohelm/x/context"
	"golang.org/x/sync/errgroup"
)

var (
	TaskPath          = cue.ParsePath("$$task")
	TaskBeforeAllPath = cue.ParsePath("$$taskBeforeAll")

	TaskPathContext = contextx.New[string](
		contextx.WithDefaultsFunc(func() string {
			return ""
		}),
	)
)

type Action = func(ctx context.Context, task Task) error

type Task interface {
	graph.Node

	DepTasks() iter.Seq[Task]

	Decode(inputs any) error
	Fill(x map[string]any) error

	Scope() Scope
	Run(ctx context.Context, action Action) error
}

type TaskUnmarshaler interface {
	UnmarshalTask(t Task) error
}

type task struct {
	graph.Node

	scope     Scope
	rank      int
	completed atomic.Uint64
	once      sync.Once
	err       error
}

func (t *task) Value() cue.Value {
	return t.scope.LookupPath(t.Path())
}

func (t *task) Scope() Scope {
	return t.scope
}

func (t *task) String() string {
	return string(cuepath.AsJSONPointer(t.Path()))
}

func (t *task) Decode(inputs any) error {
	if u, ok := inputs.(TaskUnmarshaler); ok {
		return u.UnmarshalTask(t)
	}

	if err := t.Value().Decode(inputs); err != nil {
		_, _ = fmt.Fprint(os.Stdout, t.Value().Syntax(
			cue.Concrete(false), // allow incomplete values
			cue.DisallowCycles(true),
			cue.Docs(true),
			cue.All(),
		))
		_, _ = fmt.Fprintln(os.Stdout)
		return err
	}

	return nil
}

func (t *task) Fill(x map[string]any) error {
	if err := t.scope.FillPath(t.Path(), x); err != nil {
		return err
	}
	return nil
}

func (t *task) Run(ctx context.Context, action func(ctx context.Context, task Task) error) error {
	depTasks := slices.Collect(t.DepTasks())

	if len(depTasks) > 0 {
		eg, c := errgroup.WithContext(ctx)

		for _, d := range depTasks {
			eg.Go(func() error {
				return d.Run(c, action)
			})
		}

		if err := eg.Wait(); err != nil {
			// ignore cancel
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return err
		}
	}

	t.once.Do(func() {
		t.err = action(ctx, t)
	})

	return t.err
}

func (t *task) AddDep(dep graph.Node) {
	if isTaskNode(t) {
		dep.(*task).rank++
	}

	t.Node.(graph.NodeAccessor).AddDep(dep)
}

func (t *task) DepTasks() iter.Seq[Task] {
	return func(yield func(Task) bool) {
		for d := range t.Node.Deps() {
			if isTaskNode(d) {
				if !yield(any(d).(Task)) {
					return
				}
			}
		}
	}
}

func (t *task) Shape() string {
	if isTaskNode(t.Node) {
		return "step"
	}
	return ""
}

func isTaskNode(n graph.Node) bool {
	return n.Value().LookupPath(TaskPath).Exists()
}
