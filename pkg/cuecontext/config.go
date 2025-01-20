package cuecontext

import (
	"os"
	"path/filepath"

	"github.com/octohelm/cuekit/pkg/mod/modfile"

	"cuelang.org/go/cue/load"
	"github.com/pkg/errors"

	"github.com/octohelm/cuekit/pkg/mod/modregistry"
	"github.com/octohelm/cuekit/pkg/mod/module"
)

type Config = load.Config

type OptionFunc = func(c *ConfigWithModule)

func WithRoot(dir string) OptionFunc {
	return func(c *ConfigWithModule) {
		c.Dir = dir
	}
}

func WithModule(m *module.Module) OptionFunc {
	return func(c *ConfigWithModule) {
		c.Module = m
	}
}

type ConfigWithModule struct {
	*Config
	*module.Module
}

func NewConfig(optionFns ...OptionFunc) (*ConfigWithModule, error) {
	c := &ConfigWithModule{
		Config: &Config{},
	}

	for i := range optionFns {
		optionFns[i](c)
	}

	dir, err := ResolveAbsDir(c.Dir)
	if err != nil {
		return nil, err
	}
	c.Dir = dir

	if _, err := os.Stat(c.Dir); err != nil {
		return nil, errors.Wrapf(err, "%s", c.Dir)
	}

	if c.Module == nil {
		c.Module = &module.Module{}
	}

	if c.Module.Language == nil {
		c.Module.Language = &modfile.Language{
			Version: modfile.GetCueVersion(),
		}
	}

	c.Module.SourceLoc = module.SourceLocOfOSDir(c.Dir)

	r, err := modregistry.NewRegistry(c.Module)
	if err != nil {
		return nil, err
	}

	c.Registry = r

	return c, nil
}

func ResolveAbsDir(dir string) (string, error) {
	if dir == "" || !filepath.IsAbs(dir) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		if dir != "" {
			return filepath.Join(cwd, dir), nil
		}
		return cwd, nil
	}

	return dir, nil
}
