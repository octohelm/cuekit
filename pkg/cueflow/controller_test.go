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
	. "github.com/octohelm/x/testing/v2"

	"github.com/octohelm/cuekit/pkg/cueflow"
	"github.com/octohelm/cuekit/pkg/cueflow/graph"
	"github.com/octohelm/cuekit/pkg/cueflow/runner"
	"github.com/octohelm/cuekit/pkg/cueflow/testdata/example"
	"github.com/octohelm/cuekit/pkg/cueutil"
)

var presets []byte

func init() {
	b := bytes.NewBuffer(nil)
	for t := range example.Registry.Tasks() {
		_ = t.WriteDeclTo(b)
	}

	raw, err := cueformat.Source(b.Bytes(), cueformat.Simplify())
	if err != nil {
		panic(err)
	}
	presets = raw
}

func loadWithDecl(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return slices.Concat(data, presets), nil
}

func FuzzController(f *testing.F) {
	f.Add("simple")
	f.Add("multi_arch_build")

	f.Fuzz(func(t *testing.T, taskName string) {
		t.Run(fmt.Sprintf("GIVEN tasks from %s", taskName), func(t *testing.T) {
			// 使用 MustValue 确保前置条件满足
			root := MustValue(t, func() (cue.Value, error) {
				data, err := loadWithDecl(fmt.Sprintf("./testdata/%s.cue", taskName))
				if err != nil {
					return cue.Value{}, err
				}
				return cueutil.BuildFile(data)
			})

			t.Run("WHEN init and run controller", func(t *testing.T) {
				// 1. 环境准备：改用 t.Setenv 确保并发安全
				t.Setenv("KEY", "key")

				action := runner.AsAction(example.Registry)
				ctrl := &cueflow.Controller{
					Action: func(ctx context.Context, task cueflow.Task) error {
						fmt.Println("RUN", task.Path())
						return action(ctx, task)
					},
					PrintKrokiURI: true,
				}

				ctx := logr.WithLogger(context.Background(), slog.Logger(slog.Default()))

				// 2. 核心行为与快照验证
				Then(t, "controller should run successfully and results should match snapshot",
					// 初始化
					ExpectMust(func() error {
						return ctrl.Init(root)
					}),
					// 执行
					ExpectMust(func() error {
						return ctrl.Run(ctx)
					}),
					// 快照比对
					ExpectMustValue(
						func() (Snapshot, error) {
							// 图表快照
							d2File := SnapshotFileFromRaw("g.d2", graph.ToD2Graph(graph.Collect(ctrl.Nodes())))

							// 结果快照
							ret := ctrl.LookupPath(cue.ParsePath("action.result"))
							jsonRaw, err := ret.MarshalJSON()
							if err != nil {
								return nil, err
							}
							jsonFile := SnapshotFileFromRaw("result.json", jsonRaw)

							return SnapshotOf(d2File, jsonFile), nil
						},
						MatchSnapshot(taskName),
					),
				)
			})
		})
	})
}
