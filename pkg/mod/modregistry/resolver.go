package modregistry

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/go-courier/logr"
	"github.com/octohelm/cuekit/internal/gomod"
	"github.com/octohelm/cuekit/pkg/mod/modfile"
	"github.com/octohelm/cuekit/pkg/mod/module"
)

type resolver struct {
	Root     *module.Module
	CacheDir string
	resolved sync.Map
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

func (r *resolver) ResolveLocal(ctx context.Context, path string, mv module.Version) (module.SourceLoc, error) {
	do, _ := r.resolved.LoadOrStore(fmt.Sprintf("%s@%s", mv.Path(), mv.Version()), sync.OnceValues(func() (module.SourceLoc, error) {
		src := filepath.Join(r.Root.SourceRoot(), path)

		info, err := os.Stat(src)
		if err != nil || !info.IsDir() {
			return module.SourceLoc{}, fmt.Errorf("replace source dir %s must exists: %w", src, err)
		}

		if links, ok := r.Root.Overwrites.ModLinks(mv.Path()); ok {
			for dst, link := range links {
				if err := r.syncAsLink(ctx, dst, src, link); err != nil {
					return module.SourceLoc{}, fmt.Errorf("sync link as %s failed: %w", dst, err)
				}
			}
		}

		return module.SourceLocOfOSDir(src), nil
	}))

	return do.(func() (module.SourceLoc, error))()
}

func (r *resolver) Resolve(ctx context.Context, mpath string, version string) (module.SourceLoc, error) {
	do, _ := r.resolved.LoadOrStore(fmt.Sprintf("%s@%s", mpath, version), sync.OnceValues(func() (module.SourceLoc, error) {
		info, err := r.gomodDownload(ctx, mpath, version)
		if err != nil {
			return module.SourceLoc{}, err
		}

		if links, ok := r.Root.Overwrites.ModLinks(mpath); ok {
			for dst, link := range links {
				if err := r.syncAsLink(ctx, dst, info.Dir, link); err != nil {
					return module.SourceLoc{}, fmt.Errorf("sync link as %s failed: %w", dst, err)
				}
			}
		}

		loc, err := r.convertToCueMod(ctx, mpath, info)
		if err != nil {
			return module.SourceLoc{}, err
		}
		return loc, nil
	}))

	return do.(func() (module.SourceLoc, error))()
}

func (r *resolver) gomodDownload(ctx context.Context, mpath string, version string) (*gomod.Module, error) {
	pkg := module.BasePath(mpath)
	info := gomod.Get(ctx, pkg, version, true)
	if info == nil {
		return nil, fmt.Errorf("can't found %s@%s", pkg, version)
	}
	if info.Error != "" {
		return nil, fmt.Errorf("gomod download failed: %w", errors.New(info.Error))
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

func (r *resolver) syncAsLink(ctx context.Context, dst, src string, link *modfile.Link) error {
	if dst == "" {
		return fmt.Errorf("invalid link dest: %s", dst)
	}

	outDir := filepath.Join(r.Root.SourceRoot(), dst)

	if err := os.RemoveAll(outDir); err != nil {
		return fmt.Errorf("clean %s failed: %w", outDir, err)
	}

	return fs.WalkDir(module.OSDirFS(src), link.Path, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(link.Path, path)
		if err != nil {
			return err
		}

		return copyFile(filepath.Join(src, path), filepath.Join(outDir, relPath))
	})
}

func (r *resolver) convertToCueMod(ctx context.Context, mpath string, info *gomod.Module) (module.SourceLoc, error) {
	outDir := filepath.Join(r.CacheDir, fmt.Sprintf("%s@%s", module.BasePath(mpath), info.Version))

	if err := os.RemoveAll(outDir); err != nil {
		return module.SourceLoc{}, fmt.Errorf("clean %s failed: %w", outDir, err)
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

		return copyFile(filepath.Join(info.Dir, path), filepath.Join(outDir, path))
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

	// switch to outDir
	mod.SourceLoc = module.SourceLocOfOSDir(outDir)
	if err := mod.Save(); err != nil {
		return module.SourceLoc{}, err
	}

	return mod.SourceLoc, nil
}
