package cueflow

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"slices"
	"sync"
	"sync/atomic"

	"cuelang.org/go/cue"
	cueformat "cuelang.org/go/cue/format"
	"golang.org/x/sync/errgroup"

	contextx "github.com/octohelm/x/context"

	"github.com/octohelm/cuekit/pkg/cueflow/graph"
	"github.com/octohelm/cuekit/pkg/cuepath"
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

type CueValueUnmarshaler interface {
	UnmarshalCueValue(cueValue cue.Value) error
}

type task struct {
	graph.Node

	scope     Scope
	rank      int
	completed atomic.Uint64
	once      sync.Once
	err       error
	runnable  bool
}

func (t *task) Scope() Scope {
	return t.scope
}

func (t *task) String() string {
	return string(cuepath.AsJSONPointer(t.Path()))
}

func (t *task) Decode(inputs any) error {
	if u, ok := inputs.(CueValueUnmarshaler); ok {
		return t.scope.DecodePathWith(t.Path(), func(v cue.Value) error {
			return u.UnmarshalCueValue(v)
		})
	}

	if u, ok := inputs.(TaskUnmarshaler); ok {
		return u.UnmarshalTask(t)
	}

	if err := t.scope.DecodePath(t.Path(), inputs); err != nil {
		raw, _ := cueformat.Node(t.scope.Value().Syntax(
			cue.Concrete(false), // allow incomplete values
			cue.DisallowCycles(true),
			cue.Docs(true),
			cue.All(),
		))
		fmt.Println(string(raw))
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
	if t.runnable {
		dep.(*task).rank++
	}

	t.Node.(graph.NodeAccessor).AddDep(dep)
}

func (t *task) DepTasks() iter.Seq[Task] {
	return func(yield func(Task) bool) {
		for d := range t.Node.Deps() {
			if d.(*task).runnable {
				if !yield(any(d).(Task)) {
					return
				}
			}
		}
	}
}

func (t *task) Shape() string {
	if t.runnable {
		return "step"
	}
	return ""
}

func isTaskNode(root cue.Value, n graph.Node) bool {
	return root.LookupPath(n.Path()).LookupPath(TaskPath).Exists()
}
