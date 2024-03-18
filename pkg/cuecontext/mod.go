package cuecontext

import (
	"context"
	"github.com/octohelm/cuekit/internal/cue"
	"github.com/octohelm/cuekit/pkg/mod/module"
)

func Init(ctx context.Context, moduleRoot string, mpath string) error {
	absDir, err := ResolveAbsDir(moduleRoot)
	if err != nil {
		return err
	}

	m := module.Module{}
	m.SourceLoc = module.SourceLocOfOSDir(absDir)
	m.Module = mpath

	return m.Save()
}

func Tidy(ctx context.Context, moduleRoot string) error {
	m := &module.Module{}

	c, err := NewConfig(WithRoot(moduleRoot), WithModule(m))
	if err != nil {
		return err
	}

	mf, err := cue.Tidy(ctx, m.FS, ".", c.Registry, "")
	if err != nil {
		return err
	}

	m.File = *mf

	m.Tidy()

	return m.Save()
}
