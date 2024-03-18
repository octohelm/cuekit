package cue

import (
	"context"
	"io/fs"

	"cuelang.org/go/mod/modconfig"
	"cuelang.org/go/mod/modfile"

	"github.com/octohelm/cuekit/internal/cue/internal/cuelang.org/go/internals/mod/modload"
)

func Tidy(ctx context.Context, fsys fs.FS, modRoot string, reg modconfig.Registry, version string) (*modfile.File, error) {
	return modload.Tidy(ctx, fsys, modRoot, reg, version)
}
