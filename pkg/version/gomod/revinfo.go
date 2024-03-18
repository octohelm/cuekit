package gomod

import (
	"context"
	"github.com/octohelm/cuekit/internal/gomod"
)

type RevInfo = gomod.RevInfo

func RevInfoFromDir(ctx context.Context, dir string) (*RevInfo, error) {
	return gomod.RevInfoFromDir(ctx, dir)
}
