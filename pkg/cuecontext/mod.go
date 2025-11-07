package cuecontext

import (
	"context"

	"cuelang.org/go/mod/modload"

	"github.com/octohelm/cuekit/pkg/mod/modfile"
	"github.com/octohelm/cuekit/pkg/mod/module"
)

func Init(ctx context.Context, moduleRoot string, mpath string) error {
	absDir, err := ResolveAbsDir(moduleRoot)
	if err != nil {
		return err
	}

	m := module.Module{}
	m.SourceLoc = module.SourceLocOfOSDir(absDir)
	m.Language = &modfile.Language{
		Version: modfile.GetCueVersion(),
	}
	m.Module = mpath

	return m.Save()
}

func Tidy(ctx context.Context, moduleRoot string) error {
	m := &module.Module{}

	c, err := NewConfig(WithRoot(moduleRoot), WithModule(m))
	if err != nil {
		return err
	}

	mf, err := modload.Tidy(ctx, m.FS, ".", c.Registry)
	if err != nil {
		return err
	}

	m.File = *mf
	m.Tidy()

	return m.Save()
}
