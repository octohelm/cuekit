package cfg

import (
	"context"
	"go/build"
	"io"
	"os"
	"path/filepath"
)

var ModulesEnabled bool

var (
	BuildContext     = build.Default
	BuildMod         string   // -mod flag
	BuildModExplicit bool     // whether -mod was set explicitly
	BuildModReason   string   // reason -mod was set, if set by default
	BuildN           bool     // -n flag
	BuildToolexec    []string // -toolexec flag
	BuildV           = true
	BuildX           bool   // -x flag
	ModCacheRW       bool   // -modcacherw flag
	ModFile          string // -modfile flag
	CmdName          = "get"
	GoPathError      string
)

var SumdbDir = gopathDir("pkg/sumdb")

func gopathDir(rel string) string {
	list := filepath.SplitList(BuildContext.GOPATH)
	if len(list) == 0 || list[0] == "" {
		return ""
	}
	return filepath.Join(list[0], rel)
}

func WithBuildXWriter(ctx context.Context, xLog io.Writer) context.Context {
	return context.WithValue(ctx, buildXContextKey{}, xLog)
}

type buildXContextKey struct{}

func BuildXWriter(ctx context.Context) (io.Writer, bool) {
	if !BuildX {
		return nil, false
	}
	if v := ctx.Value(buildXContextKey{}); v != nil {
		return v.(io.Writer), true
	}
	return os.Stderr, true
}
