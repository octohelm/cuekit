package cuecontext

import (
	"os"

	"path/filepath"

	"cuelang.org/go/cue/load"
	"github.com/pkg/errors"

	"github.com/octohelm/cuekit/pkg/mod/modregistry"
	"github.com/octohelm/cuekit/pkg/mod/module"
)

type Config = load.Config

type OptionFunc = func(c *ctx)

func WithRoot(dir string) OptionFunc {
	return func(c *ctx) {
		c.Dir = dir
	}
}

func WithModule(m *module.Module) OptionFunc {
	return func(c *ctx) {
		c.Module = m
	}
}

type ctx struct {
	*Config
	*module.Module
}

func NewConfig(optionFns ...OptionFunc) (*Config, error) {
	c := &ctx{
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

	c.Module.SourceLoc = module.SourceLocOfOSDir(c.Dir)

	r, err := modregistry.NewRegistry(c.Module)
	if err != nil {
		return nil, err
	}

	c.Registry = r

	return c.Config, nil
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
