package cueflow_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"slices"
	"testing"

	"cuelang.org/go/cue"
	cueformat "cuelang.org/go/cue/format"

	"github.com/octohelm/x/logr"
	"github.com/octohelm/x/logr/slog"
	"github.com/octohelm/x/testing/bdd"
	"github.com/octohelm/x/testing/snapshot"

	"github.com/octohelm/cuekit/pkg/cueflow"
	"github.com/octohelm/cuekit/pkg/cueflow/graph"
	"github.com/octohelm/cuekit/pkg/cueflow/runner"
	"github.com/octohelm/cuekit/pkg/cueflow/testdata/example"
	"github.com/octohelm/cuekit/pkg/cueutil"
)

var decls []byte

func init() {
	b := bytes.NewBuffer(nil)

	for t := range example.Registry.Tasks() {
		_ = t.WriteDeclTo(b)
	}

	decls = bdd.Must(cueformat.Source(b.Bytes(), cueformat.Simplify()))
}

func loadWithDecl(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return slices.Concat(data, decls), nil
}

func FuzzController(f *testing.F) {
	f.Add("simple")
	f.Add("multi_arch_build")

	f.Fuzz(func(t *testing.T, task string) {
		b := bdd.FromT(t)

		b.Given("tasks", func(b bdd.T) {
			root := bdd.Must(cueutil.BuildFile(bdd.Must(loadWithDecl(fmt.Sprintf("./testdata/%s.cue", task)))))

			if err := root.Err(); err != nil {
				t.Fatal("failed", err)
			}

			action := runner.AsAction(example.Registry)

			ctrl := &cueflow.Controller{
				Action: func(ctx context.Context, task cueflow.Task) error {
					fmt.Println("RUN", task.Path())

					return action(ctx, task)
				},
				PrintKrokiURI: true,
			}

			b.Then("init success",
				bdd.NoError(ctrl.Init(root)),
			)

			b.When("run tasks", func(b bdd.T) {
				_ = os.Setenv("KEY", "key")

				ctx := logr.WithLogger(b.Context(), slog.Logger(slog.Default()))

				err := ctrl.Run(ctx)
				b.Then("success", bdd.NoError(err))

				ret := ctrl.LookupPath(cue.ParsePath("action.result"))

				b.Then("got results",
					bdd.MatchSnapshot(func(s *snapshot.Snapshot) {
						s.Add("g.d2", graph.ToD2Graph(graph.Collect(ctrl.Nodes())))
						s.Add("result.json", bdd.Must(ret.MarshalJSON()))
					}, task),
				)
			})
		})
	})
}
