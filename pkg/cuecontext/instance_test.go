package cuecontext

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/octohelm/cuekit/pkg/mod/modmem"
	"github.com/octohelm/unifs/pkg/filesystem"
	testingx "github.com/octohelm/x/testing"
)

func Test(t *testing.T) {
	t.Run("should build simple", func(t *testing.T) {
		c, err := NewConfig(WithRoot("./testdata/simple"))
		testingx.Expect(t, err, testingx.Be[error](nil))

		_, err = EvalJSON(c.Config)
		testingx.Expect(t, err, testingx.Be[error](nil))
	})

	t.Run("should build with go mod", func(t *testing.T) {
		c, err := NewConfig(WithRoot("./testdata/gomod"))
		testingx.Expect(t, err, testingx.Be[error](nil))

		_, err = EvalJSON(c.Config)
		testingx.Expect(t, err, testingx.Be[error](nil))
	})

	t.Run("should build with local replace", func(t *testing.T) {
		c, err := NewConfig(WithRoot("./testdata/localreplace"))
		testingx.Expect(t, err, testingx.Be[error](nil))

		_, err = EvalJSON(c.Config)
		testingx.Expect(t, err, testingx.Be[error](nil))
	})

	t.Run("should build with mem module", func(t *testing.T) {
		c, err := NewConfig(WithRoot("./testdata/mem"))
		testingx.Expect(t, err, testingx.Be[error](nil))

		_, err = EvalJSON(c.Config)
		testingx.Expect(t, err, testingx.Be[error](nil))
	})

	t.Run("mod init && tidy", func(t *testing.T) {
		moduleRoot := "./testdata/tidy"

		_ = os.RemoveAll(filepath.Join(moduleRoot, "cache"))
		_ = os.RemoveAll(filepath.Join(moduleRoot, "cue.mod"))

		err := Init(context.Background(), moduleRoot, "ghcr.io/octothelm/tidy@v0")
		testingx.Expect(t, err, testingx.Be[error](nil))

		err = Tidy(context.Background(), moduleRoot)
		testingx.Expect(t, err, testingx.Be[error](nil))
	})
}

func init() {
	// register meme
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
