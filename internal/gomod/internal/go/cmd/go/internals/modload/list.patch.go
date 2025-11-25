package modload

import (
	"context"
	"sync"

	"github.com/octohelm/cuekit/internal/gomod/internal/go/cmd/go/internals/modinfo"
)

var listModulesMu = &sync.Mutex{}

func ListModulesOnly(state *State, ctx context.Context, args []string, mode ListMode) ([]*modinfo.ModulePublic, error) {
	listModulesMu.Lock()
	defer listModulesMu.Unlock()

	_, mods, err := listModules(state, ctx, LoadModFile(state, ctx), args, mode, nil)
	return mods, err
}
