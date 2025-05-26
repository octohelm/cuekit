package modmem

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/octohelm/cuekit/pkg/mod/modfile"
	"github.com/octohelm/cuekit/pkg/mod/module"
)

var DefaultRegistry = Registry{}

type Registry map[string]*Module

func (r Registry) Register(m *Module) {
	r[m.Module] = m
}

func (r Registry) Dump(cacheDir string) error {
	for _, m := range r {
		if err := r.dumpModuleToCache(cacheDir, m); err != nil {
			return err
		}
	}
	return nil
}

func (r Registry) Resolve(mpath string) (*module.Module, bool) {
	path := module.BasePath(mpath)

	for p, m := range r {
		if module.BasePath(p) == path {
			return &module.Module{
				File: m.File,
				Overwrites: &modfile.FileOverwrites{
					Version: m.Version,
				},
			}, true
		}
	}
	return nil, false
}

func (r Registry) dumpModuleToCache(cacheDir string, m *Module) error {
	dist := r.CacheDir(cacheDir, m.Module, m.Version)

	if err := os.RemoveAll(dist); err != nil {
		return fmt.Errorf("clean dest %s failed: %w", dist, err)
	}

	return fs.WalkDir(m.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dist, path), os.ModePerm)
		}
		src, err := m.FS.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		dest, err := os.OpenFile(filepath.Join(dist, path), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}

		defer dest.Close()

		_, err = io.Copy(dest, src)
		return err
	})
}

func (r Registry) CacheDir(cacheDir string, mpath string, version string) string {
	base, _, ok := module.SplitPathVersion(mpath)
	if !ok {
		base = mpath
	}
	return filepath.Join(cacheDir, fmt.Sprintf("%s@%s", base, version))
}
