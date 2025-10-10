package gomod

import (
	"context"
	"os"
	"testing"

	"github.com/octohelm/cuekit/internal/gomod/internal/cmd/go/internals/cfg"
	testingx "github.com/octohelm/x/testing"
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
				e := Get(ctx, p, v, true)
				t.Log(e.Path, e.Version, e.Dir, e.Sum)
				testingx.Expect(t, e.Error, testingx.Be(""))
			})
		}
	})

	t.Run("with go.mod", func(t *testing.T) {
		_ = os.RemoveAll("testdata/tmp")
		_ = os.MkdirAll("testdata/tmp", os.ModePerm)
		_ = os.Chdir("testdata/tmp")
		_ = os.WriteFile("go.mod", []byte(`module github.com/octohelm/cuekit/test
go 1.22
`), os.ModePerm)

		for p, v := range pkgs {
			t.Run("get "+p+"@"+v, func(t *testing.T) {
				e := Get(ctx, p, v, true)
				t.Log(e.Path, e.Version, e.Dir, e.Sum)
				testingx.Expect(t, e.Error, testingx.Be(""))
			})
		}
	})
}
