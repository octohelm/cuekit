package modregistry

import (
	"context"
	"fmt"
	"github.com/octohelm/cuekit/pkg/mod/modfile"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/go-courier/logr"
	"github.com/pkg/errors"

	"github.com/octohelm/cuekit/internal/gomod"
	"github.com/octohelm/cuekit/pkg/mod/module"
)

type resolver struct {
	CacheDir string
	Root     *module.Module

	resolved sync.Map
}

func (r *resolver) ResolveLocal(ctx context.Context, path string) (module.SourceLoc, error) {
	dir := filepath.Join(r.Root.SourceRoot(), path)

	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return module.SourceLoc{}, errors.Errorf("replace dir must exists dir %s", dir)
	}

	return module.SourceLocOfOSDir(dir), nil
}

func (r *resolver) Fetch(ctx context.Context, m module.Version) (module.SourceLoc, error) {
	return r.Resolve(ctx, m.Path(), m.Version())
}

var reVersionSuffixed = regexp.MustCompile("(.+)/v[\\d+]+$")

func trimVersionedSuffix(p string) string {
	if reVersionSuffixed.MatchString(p) {
		return reVersionSuffixed.FindAllStringSubmatch(p, 1)[0][1]
	}
	return p
}

func (r *resolver) ModuleVersions(ctx context.Context, mpath string) ([]string, error) {
	path := module.BasePath(mpath)

	repoRoot, err := gomod.RepoRootForImportPath(path)
	if err != nil {
		return nil, nil
	}

	if repoRoot != trimVersionedSuffix(path) {
		return nil, nil
	}

	versions, _ := gomod.ListVersion(ctx, path)
	if len(versions) > 0 {
		validVersions := make([]string, 0, len(versions))

		for _, version := range versions {
			v, err := module.NewVersion(path, version)
			if err == nil {
				r.Root.Overwrites.AddDep(v.Path(), &modfile.DepOverwrite{
					Path:    path,
					Version: v.Version(),
				})
				validVersions = append(validVersions, v.Version())
			} else {
				panic(err)
			}
		}

		return validVersions, nil

	}
	return versions, nil
}

func (r *resolver) Resolve(ctx context.Context, mpath string, version string) (module.SourceLoc, error) {
	do, _ := r.resolved.LoadOrStore(fmt.Sprintf("%s@%s", mpath, version), sync.OnceValue(func() any {
		info, err := r.gomodDownload(ctx, mpath, version)
		if err != nil {
			return err
		}
		loc, err := r.convertToCueMod(ctx, mpath, info)
		if err != nil {
			return err
		}
		return loc
	}))
	switch x := do.(func() any)().(type) {
	case error:
		return module.SourceLoc{}, x
	default:
		return x.(module.SourceLoc), nil
	}
}

func (r *resolver) gomodDownload(ctx context.Context, mpath string, version string) (*gomod.Module, error) {
	pkg := module.BasePath(mpath)
	info := gomod.Get(ctx, pkg, version, true)
	if info == nil {
		return nil, fmt.Errorf("can't found %s@%s", pkg, version)
	}
	if info.Error != "" {
		return nil, errors.Wrap(errors.New(info.Error), "gomod download failed")
	}

	if info.Dir == "" {
		logr.FromContext(ctx).Debug(fmt.Sprintf("get %s@%s", pkg, version))
		gomod.Download(ctx, info)
		if info.Error != "" {
			return nil, errors.New(info.Error)
		}
	}
	return info, nil
}

func (r *resolver) convertToCueMod(ctx context.Context, mpath string, info *gomod.Module) (module.SourceLoc, error) {
	dist := filepath.Join(r.CacheDir, fmt.Sprintf("%s@%s", module.BasePath(mpath), info.Version))

	if err := os.RemoveAll(dist); err != nil {
		return module.SourceLoc{}, errors.Wrapf(err, "clean dest failed: %s", dist)
	}

	if err := fs.WalkDir(module.OSDirFS(info.Dir), ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if strings.Contains(path, "cue.mod") {
				return fs.SkipDir
			}
			return nil
		}

		// only cp cue files
		if !strings.HasSuffix(path, ".cue") {
			return nil
		}

		return copyFile(filepath.Join(info.Dir, path), filepath.Join(dist, path))
	}); err != nil {
		return module.SourceLoc{}, err
	}

	mod := &module.Module{
		SourceLoc: module.SourceLocOfOSDir(info.Dir),
	}

	// could empty
	_ = mod.Load(false)

	mod.Module = mpath
	mod.Overwrites.Path = info.Path
	mod.Overwrites.Version = info.Version

	// switch to dist
	mod.SourceLoc = module.SourceLocOfOSDir(dist)
	if err := mod.Save(); err != nil {
		return module.SourceLoc{}, err
	}

	return mod.SourceLoc, nil
}
