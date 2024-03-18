package modload

import (
	"context"

	"golang.org/x/mod/module"

	"github.com/octohelm/cuekit/internal/gomod/internal/cmd/go/internals/modinfo"
)

func ListModulesOnly(ctx context.Context, args []string, mode ListMode) ([]*modinfo.ModulePublic, error) {
	var reuse map[module.Version]*modinfo.ModulePublic
	_, mods, err := listModules(ctx, LoadModFile(ctx), args, mode, reuse)
	return mods, err
}
