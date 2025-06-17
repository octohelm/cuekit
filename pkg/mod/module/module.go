package module

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	cuemodfile "cuelang.org/go/mod/modfile"
	"cuelang.org/go/mod/module"
	"github.com/octohelm/cuekit/pkg/mod/modfile"
)

const (
	fileModule           = "cue.mod/module.cue"
	fileModuleOverwrites = "cue.mod/module_overwrites.cue"
)

type OSRootFS = module.OSRootFS

type Module struct {
	module.SourceLoc
	modfile.File

	Overwrites *modfile.FileOverwrites
}

func (m *Module) SourceRoot() string {
	if osRoot, ok := m.FS.(OSRootFS); ok {
		return filepath.Join(osRoot.OSRoot(), m.Dir)
	}

	return fmt.Sprintf("mem:%s", m.Dir)
}

func (m *Module) SetDefaults() {
	if m.Language == nil {
		m.Language = &modfile.Language{
			Version: modfile.GetCueVersion(),
		}
	}
}

func (m *Module) GetDepOverwrite(mpath string) (*modfile.DepOverwrite, bool) {
	v, ok := m.Overwrites.GetDep(mpath)
	if ok {
		return v, ok
	}
	return nil, false
}

func (m *Module) Load(strict bool) error {
	if m.Overwrites == nil {
		m.Overwrites = &modfile.FileOverwrites{}
	}

	if err := m.loadModule(strict); err != nil {
		return err
	}

	// optional load
	if data, err := fs.ReadFile(m.FS, filepath.Join(m.Dir, fileModuleOverwrites)); err == nil {
		if err := m.Overwrites.Load(data); err != nil {
			return err
		}
	}

	return nil
}

func (m *Module) loadModule(strict bool) error {
	data, err := fs.ReadFile(m.FS, filepath.Join(m.Dir, fileModule))
	if err != nil {
		if os.IsNotExist(err) && !strict {
			return nil
		}
		return err
	}

	parse := cuemodfile.Parse
	if !strict {
		parse = cuemodfile.ParseNonStrict
	}

	mf, err := parse(data, "")
	if err != nil {
		return fmt.Errorf("cannot parse module file: %s, %v", m.SourceRoot(), err)
	}

	m.File.Module = mf.Module
	m.File.Deps = mf.Deps
	m.File.Language = mf.Language

	return nil
}

func (mm *Module) Save() error {
	m := &Module{
		SourceLoc:  mm.SourceLoc,
		File:       mm.File,
		Overwrites: mm.Overwrites,
	}

	m.SetDefaults()

	if !m.Overwrites.IsZero() {
		for mpath, d := range m.Overwrites.Deps {
			if d.IsLocalReplacement() {
				m.addDep(mpath, &cuemodfile.Dep{
					Version: "v0.0.0",
				})
				continue
			}

			if d.Version != "" {
				m.addDep(mpath, &cuemodfile.Dep{
					Version: d.Version,
				})
			}
		}
		data, err := m.Overwrites.Format()
		if err != nil {
			return nil
		}
		if err := putFileContents(filepath.Join(m.SourceRoot(), fileModuleOverwrites), bytes.NewBuffer(data)); err != nil {
			return err
		}
	}

	data, err := m.File.Format()
	if err != nil {
		return err
	}
	if err := putFileContents(filepath.Join(m.SourceRoot(), fileModule), bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}

func (m *Module) addDep(mpath string, dep *cuemodfile.Dep) {
	if m.File.Deps == nil {
		m.File.Deps = map[string]*cuemodfile.Dep{}
	}

	m.File.Deps[mpath] = dep
}

func (m *Module) Tidy() {
	if len(m.File.Deps) > 0 {
		for mpath := range m.Overwrites.Deps {
			if _, ok := m.File.Deps[mpath]; !ok {
				// remove unused dep overwrite
				delete(m.Overwrites.Deps, mpath)
			}
		}
	}
}

func putFileContents(filename string, r io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	df, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, r)
	return nil
}
