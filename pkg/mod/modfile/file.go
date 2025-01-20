package modfile

import (
	"runtime/debug"
	"sync"

	cuemodfile "cuelang.org/go/mod/modfile"
)

var GetCueVersion = sync.OnceValue(func() string {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range bi.Deps {
			if dep.Path == "cuelang.org/go" {
				return dep.Version
			}
		}
	}
	return "v0.11.1"
})

type (
	File     = cuemodfile.File
	Dep      = cuemodfile.Dep
	Language = cuemodfile.Language
)

type OldFile struct{}
