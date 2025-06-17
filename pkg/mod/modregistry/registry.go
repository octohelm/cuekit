package modregistry

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/mod/modconfig"
	"github.com/octohelm/cuekit/pkg/mod/modfile"
	"github.com/octohelm/cuekit/pkg/mod/modmem"
	"github.com/octohelm/cuekit/pkg/mod/module"
)

func NewRegistry(m *module.Module) (modconfig.Registry, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	if err := m.Load(true); err != nil {
		return nil, err
	}

	r := &resolver{
		CacheDir: getEnv("CUE_CACHE_DIR", filepath.Join(userCacheDir, "cue")),
	}
	if err := r.Init(m); err != nil {
		return nil, err
	}

	std, err := modconfig.NewRegistry(&modconfig.Config{
		Env: []string{
			fmt.Sprintf("CUE_REGISTRY=%s", getEnv("CUE_REGISTRY", "ghcr.io")),
			fmt.Sprintf("CUE_CACHE_DIR=%s", r.CacheDir),
		},
	})
	if err != nil {
		return nil, err
	}

	if err := modmem.DefaultRegistry.Dump(r.CacheDir); err != nil {
		return nil, err
	}

	return &registry{
		mem:   modmem.DefaultRegistry,
		local: r,
		std:   std,
	}, nil
}

type registry struct {
	mem modmem.Registry
	std modconfig.Registry

	local *resolver
}

func (r *registry) Fetch(ctx context.Context, mv module.Version) (loc module.SourceLoc, err error) {
	defer func() {
		if err == nil {
			mod := &module.Module{
				SourceLoc: loc,
			}

			if err := mod.Load(true); err == nil {
				if version := mod.Overwrites.Version; version != "" {
					r.local.AddOverwrite(mod.File.Module, &modfile.DepOverwrite{
						Path:    mod.Overwrites.Path,
						Version: version,
					})
				}
			}

			_ = r.local.Commit()
		}
	}()

	if depOverwrite, ok := r.local.GetDepOverwrite(mv.Path()); ok {
		if depOverwrite.IsLocalReplacement() {
			return r.local.ResolveLocal(ctx, depOverwrite.Path, mv)
		}

		if depOverwrite.Path != "" {
			v, err := module.NewVersion(depOverwrite.Path, depOverwrite.Version)
			if err == nil {
				mv = v
			}
		}

		return r.local.Resolve(ctx, mv.Path(), depOverwrite.Version)
	}

	if m, ok := r.mem.Resolve(mv.Path()); ok {
		return module.SourceLocOfOSDir(r.mem.CacheDir(r.local.CacheDir, m.File.Module, mv.Version())), nil
	}

	sl, err := r.std.Fetch(ctx, mv)
	if err != nil {
		if r.isNotExistsOfCueRegistry(err) {
			resp, err := r.local.Fetch(ctx, mv)
			return resp, err
		}
		return module.SourceLoc{}, err
	}
	return sl, nil
}

func (r *registry) isNotExistsOfCueRegistry(err error) bool {
	if err != nil {
		errMsg := err.Error()
		return strings.Contains(errMsg, "HTTP response 403") || strings.Contains(errMsg, "HTTP response 400")
	}
	return false
}

func (r *registry) Requirements(ctx context.Context, mv module.Version) ([]module.Version, error) {
	if m, ok := r.mem.Resolve(mv.Path()); ok {
		m.SourceLoc = module.SourceLocOfOSDir(
			r.mem.CacheDir(r.local.CacheDir, m.File.Module, m.Overwrites.Version),
		)

		if err := m.Load(true); err != nil {
			return nil, err
		}
		return m.File.DepVersions(), nil
	}

	if depOverwrite, ok := r.local.GetDepOverwrite(mv.Path()); ok {
		if depOverwrite.IsLocalReplacement() {
			s, err := r.local.ResolveLocal(ctx, depOverwrite.Path, mv)
			if err != nil {
				return nil, err
			}
			m := &module.Module{SourceLoc: s}
			if err := m.Load(false); err != nil {
				return nil, err
			}
			return m.File.DepVersions(), nil
		}

		if depOverwrite.Path != "" {
			v, err := module.NewVersion(depOverwrite.Path, depOverwrite.Version)
			if err == nil {
				mv = v
			}
		}

		s, err := r.local.Fetch(ctx, mv)
		if err != nil {
			return nil, err
		}
		m := &module.Module{
			SourceLoc: s,
		}
		if err := m.Load(false); err != nil {
			return nil, err
		}
		return m.File.DepVersions(), nil
	}

	return r.std.Requirements(ctx, mv)
}

func (r *registry) ModuleVersions(ctx context.Context, mpath string) ([]string, error) {
	if m, ok := r.mem.Resolve(mpath); ok {
		return []string{m.Overwrites.Version}, nil
	}

	if depOverwrite, ok := r.local.GetDepOverwrite(mpath); ok {
		if depOverwrite.IsLocalReplacement() {
			return []string{"v0.0.0"}, nil
		}

		if depOverwrite.Path != "" {
			mpath = depOverwrite.Path
		}
	}

	versions, _ := r.local.ModuleVersions(ctx, mpath)
	if len(versions) > 0 {
		return versions, nil
	}

	versions2, _ := r.std.ModuleVersions(ctx, mpath)
	return versions2, nil
}
