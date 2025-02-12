package modload

import (
	"context"
	"github.com/octohelm/cuekit/internal/gomod/internal/cmd/go/internals/modinfo"
)

func ListModulesOnly(ctx context.Context, args []string, mode ListMode) ([]*modinfo.ModulePublic, error) {
	_, mods, err := listModules(ctx, LoadModFile(ctx), args, mode, nil)
	return mods, err
}
