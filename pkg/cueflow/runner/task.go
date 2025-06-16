package runner

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"reflect"
	"strings"

	"cuelang.org/go/cue"
	"github.com/fatih/color"
	"github.com/go-courier/logr"
	"github.com/octohelm/cuekit/pkg/cueconvert"
	"github.com/octohelm/cuekit/pkg/cueflow"
	"github.com/octohelm/cuekit/pkg/cuepath"
	"github.com/octohelm/cuekit/pkg/cueutil"
)

type Task interface {
	NewRunner(t cueflow.Task) (TaskRunner, error)
	WriteDeclTo(w io.Writer) error
}

type TaskDoer interface {
	Do(ctx context.Context) error
}

type WithScopeName interface {
	ScopeName(ctx context.Context) (string, error)
}

type TaskRunner interface {
	Path() cue.Path
	Underlying() any

	Run(ctx context.Context) error
}

type named struct {
	tpe  reflect.Type
	decl *cueconvert.Decl

	flowStruct        bool
	flowTask          bool
	flowTaskBeforeAll bool
	flowControl       bool

	outputFields map[string][]int
}

func (t *named) WriteDeclTo(w io.Writer) error {
	if len(t.decl.Source) == 0 {
		return nil
	}

	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}

	if t.flowTaskBeforeAll {
		if _, err := fmt.Fprintf(w, `%s: $$taskBeforeAll: true
`,
			t.decl.Name,
		); err != nil {
			return err
		}
	}

	if t.flowControl {
		if _, err := fmt.Fprintf(w, `%s: $$control: name: %q
`,
			t.decl.Name,
			strings.ToLower(strings.Trim(t.decl.Name, "#")),
		); err != nil {
			return err
		}
	} else if t.flowTask {
		if _, err := fmt.Fprintf(w, `%s: $$task: true
`,
			t.decl.Name,
		); err != nil {
			return err
		}
	}

	if t.flowStruct && t.decl.Source[0] == '{' {
		if _, err := fmt.Fprintf(w, `%s: $$type: name: %q
`,
			t.decl.Name,
			t.FullName(),
		); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, `%s: %s
`,
		t.decl.Name,
		t.decl.Source,
	); err != nil {
		return err
	}

	return nil
}

func (t *named) FullName() string {
	if t.decl.PkgPath == "" {
		return t.decl.Name
	}
	return fmt.Sprintf("%s.%s", t.decl.PkgPath, t.decl.Name)
}

func (t *named) NewRunner(task cueflow.Task) (TaskRunner, error) {
	return &taskRunner{
		task:         task,
		rv:           reflect.New(t.tpe),
		outputFields: maps.Clone(t.outputFields),
	}, nil
}

type taskRunner struct {
	rv           reflect.Value
	task         cueflow.Task
	outputFields cueconvert.OutputFields
}

func (t *taskRunner) Underlying() any {
	return t.rv.Interface()
}

func (t *taskRunner) Path() cue.Path {
	return t.task.Path()
}

func (t *taskRunner) Task() cueflow.Task {
	return t.task
}

func (t *taskRunner) Run(ctx context.Context) (err error) {
	v := t.rv.Interface()

	doer := v.(TaskDoer)

	markSkip := false

	if err := t.task.Scope().DecodePathWith(t.Path(), func(cv cue.Value) error {
		dep := cv.LookupPath(PathDep)

		if ctrl := dep.LookupPath(PathControl); ctrl.Exists() {
			ctrlType, _ := ctrl.String()
			switch ctrlType {
			case "skip":
				needSkip, _ := dep.LookupPath(cue.ParsePath("when")).Bool()
				if needSkip {
					if _, ok := doer.(cueflow.TaskFeedback); ok {
						markSkip = true
					}
					return nil
				}
			}
		}

		return nil
	}); err != nil {
		return err
	}

	if markSkip {
		return t.task.Fill(map[string]any{"$ok": false})
	}

	ctx = cueflow.TaskPathContext.Inject(ctx, t.task.Path().String())

	ctx, l := logr.FromContext(ctx).Start(ctx, fmt.Sprintf("%s %s", t.task.Path().String(), color.WhiteString(t.task.Name())))
	defer l.End()

	if err := t.task.Decode(doer); err != nil {
		return fmt.Errorf("%s: decode into %T failed: %w", cuepath.AsJSONPointer(t.task.Path()), doer, err)
	}

	if n, ok := doer.(WithScopeName); ok {
		scopeName, err := n.ScopeName(ctx)
		if err != nil {
			return err
		}
		l = l.WithValues(slog.String("$scope", scopeName))
	}

	isCheckpoint := false
	if checkpoint, ok := doer.(cueflow.Checkpoint); ok {
		isCheckpoint = checkpoint.AsCheckpoint()
	}

	if !isCheckpoint {
		if err := doer.Do(ctx); err != nil {
			return fmt.Errorf("task do failed: %w", err)
		}
	}

	// done task before resolve output values
	if taskFeedback, ok := doer.(cueflow.TaskFeedback); ok {
		taskFeedback.Done(nil)
	}

	values := t.outputFields.OutputValues(t.rv)

	l.WithValues(slog.Any("values", cueutil.AsLogValue(values))).Debug("done")

	if err := t.task.Fill(values); err != nil {
		return fmt.Errorf("fill result failed: %w", err)
	}

	return nil
}
