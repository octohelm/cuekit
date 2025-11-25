package cuecontext

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/octohelm/unifs/pkg/filesystem"
	. "github.com/octohelm/x/testing/v2"

	"github.com/octohelm/cuekit/pkg/mod/modmem"
)

func Test(t *testing.T) {
	t.Run("should build simple", func(t *testing.T) {
		c := MustValue(t, func() (*ConfigWithModule, error) {
			return NewConfig(WithRoot("./testdata/simple"))
		})

		Then(t, "config should be created and evaluated",
			ExpectMust(func() error {
				_, err := EvalJSON(c.Config)
				return err
			}),
		)
	})

	t.Run("should build with go mod", func(t *testing.T) {
		Then(t, "build with go mod should succeed",
			ExpectMust(func() error {
				c, err := NewConfig(WithRoot("./testdata/gomod"))
				if err != nil {
					return err
				}
				_, err = EvalJSON(c.Config)
				return err
			}),
		)
	})

	t.Run("should build with local replace", func(t *testing.T) {
		Then(t, "build with local replace should succeed",
			ExpectMust(func() error {
				c, err := NewConfig(WithRoot("./testdata/localreplace"))
				if err != nil {
					return err
				}
				_, err = EvalJSON(c.Config)
				return err
			}),
		)
	})

	t.Run("should build with mem module", func(t *testing.T) {
		Then(t, "build with mem module should succeed",
			ExpectMust(func() error {
				c, err := NewConfig(WithRoot("./testdata/mem"))
				if err != nil {
					return err
				}
				_, err = EvalJSON(c.Config)
				return err
			}),
		)
	})

	t.Run("mod init && tidy", func(t *testing.T) {
		moduleRoot := "./testdata/tidy"

		_ = os.RemoveAll(filepath.Join(moduleRoot, "cache"))
		_ = os.RemoveAll(filepath.Join(moduleRoot, "cue.mod"))

		Then(t, "init and tidy should success",
			ExpectMust(func() error {
				return Init(context.Background(), moduleRoot, "ghcr.io/octothelm/tidy@v0")
			}),
			ExpectMust(func() error {
				return Tidy(context.Background(), moduleRoot)
			}),
		)
	})
}

func init() {
	m, _ := modmem.NewModule("mem.octothelm.tech", "v0.0.0", func(ctx context.Context, fsys filesystem.FileSystem) error {
		_ = filesystem.MkdirAll(ctx, fsys, "x")
		_ = filesystem.Write(ctx, fsys, "x/x.cue", []byte(`
package x

#Version: "mem-devel" 
`))
		return nil
	})
	modmem.DefaultRegistry.Register(m)
}
