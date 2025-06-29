package gomod

import (
	"context"
	"fmt"
	"strings"
	_ "unsafe"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"

	"github.com/octohelm/cuekit/pkg/version"

	"github.com/octohelm/cuekit/internal/gomod/internal/cmd/go/internals/modfetch"
	"github.com/octohelm/cuekit/internal/gomod/internal/cmd/go/internals/modfetch/codehost"
	"github.com/octohelm/cuekit/internal/gomod/internal/cmd/go/internals/vcs"
)

//go:linkname newCodeRepo github.com/octohelm/cuekit/internal/gomod/internal/cmd/go/internals/modfetch.newCodeRepo
func newCodeRepo(code codehost.Repo, codeRoot, path string) (modfetch.Repo, error)

//go:linkname lookupCodeRepo github.com/octohelm/cuekit/internal/gomod/internal/cmd/go/internals/modfetch.lookupCodeRepo
func lookupCodeRepo(ctx context.Context, rr *vcs.RepoRoot) (codehost.Repo, error)

func finalLookupCodeRepo(ctx context.Context, rr *vcs.RepoRoot, localOk bool) (codehost.Repo, error) {
	if strings.ToLower(rr.VCS.Name) == "git" && localOk {
		return codehost.NewRepo(ctx, "git", rr.Root, true)
	}
	return lookupCodeRepo(ctx, rr)
}

type RevInfo = modfetch.RevInfo

func RevInfoFromDir(ctx context.Context, dir string) (*RevInfo, error) {
	rootDir, c, err := vcs.FromDir(dir, "", true)
	if err != nil {
		return nil, err
	}

	repo, err := c.RemoteRepo(c, rootDir)
	if err != nil {
		return nil, fmt.Errorf("resolve remote repo failed: %w", err)
	}

	head, err := c.Status(c, rootDir)
	if err != nil {
		return nil, fmt.Errorf("stat faield: %w", err)
	}

	rr := &vcs.RepoRoot{}
	rr.VCS = c
	rr.Root = rootDir
	rr.Repo = repo

	code, err := finalLookupCodeRepo(ctx, rr, true)
	if err != nil {
		return nil, err
	}

	importPath := rr.Root

	data, err := code.ReadFile(ctx, head.Revision, "go.mod", -1)
	if err == nil {
		f, err := modfile.ParseLax("go.mod", data, nil)
		if err != nil {
			return nil, err
		}

		// <import_path>/v2
		_, pathMajor, ok := module.SplitPathVersion(f.Module.Mod.Path)
		if ok && pathMajor != "" {
			importPath += pathMajor
		}
	}

	r, err := newCodeRepo(code, rr.Root, importPath)
	if err != nil {
		return nil, err
	}

	info, err := r.Stat(ctx, head.Revision)
	if err != nil {
		return nil, fmt.Errorf("stat faield: %w", err)
	}

	info.Version = version.Convert(info.Version, info.Time, info.Short, head.Uncommitted)

	return info, nil
}
