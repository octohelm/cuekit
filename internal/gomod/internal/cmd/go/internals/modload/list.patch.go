package modload

import (
	"context"
	"sync"

	"github.com/octohelm/cuekit/internal/gomod/internal/cmd/go/internals/modinfo"
)

var listModulesMu = &sync.Mutex{}

func ListModulesOnly(ctx context.Context, args []string, mode ListMode) ([]*modinfo.ModulePublic, error) {
	listModulesMu.Lock()
	defer listModulesMu.Unlock()

	_, mods, err := listModules(ctx, LoadModFile(ctx), args, mode, nil)
	return mods, err
}
