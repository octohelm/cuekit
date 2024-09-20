package modmem

import (
	"context"

	"cuelang.org/go/mod/module"
	"github.com/octohelm/cuekit/pkg/mod/modfile"
	"github.com/octohelm/unifs/pkg/filesystem"
)

func NewModule(path string, version string, commitSource func(ctx context.Context, fsys filesystem.FileSystem) error) (*Module, error) {
	v, err := module.NewVersion(path, version)
	if err != nil {
		return nil, err
	}

	m := &Module{}
	m.Module = v.Path()
	m.Version = v.Version()
	m.Language = &modfile.Language{
		Version: modfile.CueVersion,
	}
	m.Dir = "."

	fsys := filesystem.NewMemFS()

	ctx := context.Background()

	if err := commitSource(ctx, fsys); err != nil {
		return nil, err
	}

	if err := filesystem.MkdirAll(ctx, fsys, "cue.mod"); err != nil {
		return nil, err
	}

	data, err := m.File.Format()
	if err != nil {
		return nil, err
	}

	if err := filesystem.Write(context.Background(), fsys, "cue.mod/module.cue", data); err != nil {
		return nil, err
	}

	m.FS = filesystem.AsReadDirFS(fsys)

	return m, nil
}

type Module struct {
	modfile.File
	Version string
	module.SourceLoc
}
