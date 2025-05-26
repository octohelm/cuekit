package cueflow

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"maps"

	"cuelang.org/go/cue"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/octohelm/cuekit/pkg/cueflow/graph"
	"github.com/octohelm/cuekit/pkg/cuepath"
	"golang.org/x/sync/errgroup"
)

func RunSubTasks(ctx context.Context, scope Scope, isPrefix func(p cue.Path) (bool, cue.Path)) error {
	c := scope.(*Controller)

	cc := &Controller{
		Action:        c.Action,
		PrintKrokiURI: c.PrintKrokiURI,
		IsPrefix:      isPrefix,
		scope:         c.scope,
	}

	if err := cc.Init(scope.Value()); err != nil {
		return err
	}

	if err := cc.Run(ctx); err != nil {
		return err
	}

	return nil
}

type Controller struct {
	graph.Graph

	Action   Action
	IsPrefix func(p cue.Path) (bool, cue.Path)

	PrintKrokiURI bool

	tasks map[jsontext.Pointer]*task

	*scope
}

func (c *Controller) Init(v cue.Value) error {
	if c.scope == nil {
		c.scope = &scope{root: v}
	}

	r := &graph.Resolver{
		CreateNode: func(n graph.Node) graph.Node {
			return &task{
				Node:  n,
				scope: c,
			}
		},
		IsPrefix: c.IsPrefix,
	}

	if err := r.Init(v); err != nil {
		return err
	}

	c.Graph = r
	c.tasks = map[jsontext.Pointer]*task{}

	for n := range c.Nodes() {
		if isTaskNode(n) {
			c.tasks[cuepath.AsJSONPointer(n.Path())] = n.(*task)
		}
	}

	return nil
}

func (c *Controller) Run(pctx context.Context) error {
	eg, ctx := errgroup.WithContext(pctx)

	for t := range c.Tasks() {
		// only trigger the final task
		if t.(*task).rank == 0 {
			eg.Go(func() error {
				return t.Run(ctx, func(ctx context.Context, task Task) error {
					return c.Action(ctx, task)
				})
			})
		}
	}

	if c.PrintKrokiURI {
		krokiURI, err := graph.ToKrokiURI(graph.ToD2Graph(graph.Collect(c.Nodes())))
		if err != nil {
			return err
		}
		fmt.Println(krokiURI)
	}

	if err := eg.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}

	return nil
}

func (c *Controller) RunMatched(pctx context.Context, match func(t Task) bool) error {
	eg, ctx := errgroup.WithContext(pctx)

	nodes := map[jsontext.Pointer]graph.Node{}

	var collectToPrint func(n graph.Node)

	if c.PrintKrokiURI {
		collectToPrint = func(n graph.Node) {
			nodes[cuepath.AsJSONPointer(n.Path())] = n

			for d := range n.Deps() {
				collectToPrint(d)
			}
		}
	}

	for t := range c.Tasks() {
		// only trigger the matched task
		if match(t) {
			if collectToPrint != nil {
				collectToPrint(t)
			}

			eg.Go(func() error {
				return t.Run(ctx, func(ctx context.Context, task Task) error {
					return c.Action(ctx, task)
				})
			})
		}
	}

	if c.PrintKrokiURI {
		krokiURI, err := graph.ToKrokiURI(graph.ToD2Graph(graph.Collect(maps.Values(nodes))))
		if err != nil {
			return err
		}
		fmt.Println(krokiURI)
	}

	if err := eg.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}

	return nil
}

func (c *Controller) Tasks() iter.Seq[Task] {
	return func(yield func(Task) bool) {
		for k := range maps.Keys(c.tasks) {
			if !yield(c.tasks[k]) {
				return
			}
		}
	}
}

func (c *Controller) LookupTask(p cue.Path) (Task, bool) {
	if c.tasks == nil {
		return nil, false
	}
	if t, ok := c.tasks[cuepath.AsJSONPointer(p)]; ok {
		return t, true
	}
	return nil, false
}
