package gomod

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/octohelm/x/cmp"
	. "github.com/octohelm/x/testing/v2"

	"github.com/octohelm/cuekit/internal/gomod/internal/cfg"
)

func TestDownload(t *testing.T) {
	_ = os.RemoveAll(cfg.GOMODCACHE)

	ctx := context.Background()

	pkgs := map[string]string{
		"github.com/octohelm/crkit":   "",
		"github.com/octohelm/kubepkg": "latest",
		"k8s.io/api":                  "v0.24.0",
	}

	t.Run("without go.mod", func(t *testing.T) {
		for p, v := range pkgs {
			t.Run("get "+p+"@"+v, func(t *testing.T) {
				Then(t, "should get package info",
					ExpectMustValue(func() (*Module, error) {
						e := Get(ctx, NewState(), p, v, true)
						if e.Error != "" {
							return nil, errors.New(e.Error)
						}
						t.Logf("Path: %s, Version: %s, Dir: %s, Sum: %s", e.Path, e.Version, e.Dir, e.Sum)
						return e, nil
					},
						Be(cmp.NotZero[*Module]()),
					),
				)
			})
		}
	})

	t.Run("with go.mod", func(t *testing.T) {
		tmpDir := t.TempDir()

		cwd, _ := os.Getwd()
		defer func() { _ = os.Chdir(cwd) }()

		_ = os.Chdir(tmpDir)
		_ = os.WriteFile("go.mod", []byte(`
module github.com/octohelm/cuekit/test

go 1.26
`), 0644)

		for p, v := range pkgs {
			t.Run("get "+p+"@"+v, func(t *testing.T) {
				Then(t, "should get package info with go.mod context",
					ExpectMustValue(func() (*Module, error) {
						e := Get(ctx, NewState(), p, v, true)
						if e.Error != "" {
							return nil, errors.New(e.Error)
						}
						return e, nil
					},
						Be(cmp.NotZero[*Module]()),
					),
				)
			})
		}
	})
}
