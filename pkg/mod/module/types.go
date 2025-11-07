package module

import (
	"cuelang.org/go/mod/module"
)

type Version = module.Version

func NewVersion(path string, version string) (Version, error) {
	return module.NewVersion(path, version)
}

type SourceLoc = module.SourceLoc
