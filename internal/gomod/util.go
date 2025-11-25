package gomod

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"sync"

	"golang.org/x/mod/module"

	"github.com/octohelm/cuekit/internal/gomod/internal/cfg"
	"github.com/octohelm/cuekit/internal/gomod/internal/go/cmd/go/internals/gover"
	"github.com/octohelm/cuekit/internal/gomod/internal/go/cmd/go/internals/modfetch"
	"github.com/octohelm/cuekit/internal/gomod/internal/go/cmd/go/internals/modload"
	"github.com/octohelm/cuekit/internal/gomod/internal/go/cmd/go/internals/vcs"
	"github.com/octohelm/cuekit/internal/gomod/internal/go/cmd/go/internals/web"
)

type Module struct {
	Path    string
	Version string
	Error   string
	Dir     string
	Sum     string
}

func init() {
	// hack to ignore patch version check
	gover.Startup.GOTOOLCHAIN = "auto"
	gover.TestVersion = "go1." + strconv.Itoa(math.MaxInt16) + "." + strconv.Itoa(math.MaxInt16)
}

type RepoRoot = vcs.RepoRoot

func RepoRootForImportPath(importPath string) (*RepoRoot, error) {
	r, err := vcs.RepoRootForImportPath(importPath, vcs.IgnoreMod, web.DefaultSecurity)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errors.New("no repo root for import path " + importPath)
	}
	return r, nil
}

var goSettingOnce = &sync.Once{}

type State = modload.State

func NewState() *State {
	state := modload.NewState()
	modload.Init(state)
	state.ForceUseModules = true

	return state
}

// Get Module
func Get(ctx context.Context, state *State, path string, version string, verbose bool) *Module {
	goSettingOnce.Do(func() {
		cfg.BuildX = verbose
	})

	mv := module.Version{Path: path, Version: version}
	p, err := modfetch.DownloadDir(ctx, mv)
	if err == nil {
		// found in cache
		return &Module{
			Path:    mv.Path,
			Version: mv.Version,
			Dir:     p,
			Sum:     modfetch.Sum(ctx, mv),
		}
	}

	if version == "" {
		version = "latest"
	}

	requestPath := path + "@" + version

	found, err := modload.ListModulesOnly(state, ctx, []string{requestPath}, modload.ListVersions)
	if err != nil {
		panic(fmt.Errorf("list %s failed: %w", requestPath, err))
	}
	if len(found) > 0 {
		info := found[0]

		m := &Module{
			Path:    info.Path,
			Version: info.Version,
		}

		if info.Error != nil {
			m.Error = info.Error.Err
		} else {
			m.Dir = info.Dir
			m.Sum = modfetch.Sum(ctx, module.Version{Path: m.Path, Version: m.Version})
		}
		return m
	}
	return nil
}

// Download Module
func Download(ctx context.Context, state *State, m *Module) {
	mv := module.Version{Path: m.Path, Version: m.Version}
	dir, err := modfetch.DownloadDir(ctx, mv)
	if err == nil {
		// found in cache
		m.Dir = dir
		m.Sum = modfetch.Sum(ctx, module.Version{Path: m.Path, Version: m.Version})
		return
	}

	dir, err = state.Fetcher().Download(ctx, mv)
	if err != nil {
		m.Error = err.Error()
		return
	}
	m.Dir = dir
	m.Sum = modfetch.Sum(ctx, module.Version{Path: m.Path, Version: m.Version})
}

func ListVersion(ctx context.Context, state *State, path string) ([]string, error) {
	found, err := modload.ListModulesOnly(state, ctx, []string{path}, modload.ListVersions)
	if err != nil {
		return nil, err
	}
	if len(found) > 0 {
		info := found[0]
		if len(info.Versions) > 0 {
			return info.Versions, nil
		}

		m := Get(ctx, state, info.Path, "latest", false)
		if m.Error != "" {
			return nil, errors.New(m.Error)
		}
		return []string{m.Version}, nil
	}
	return nil, errors.New("no versions")
}
