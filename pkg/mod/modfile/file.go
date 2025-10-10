package modfile

import (
	"sync"

	cuemodfile "cuelang.org/go/mod/modfile"
)

var GetCueVersion = sync.OnceValue(func() string {
	//bi, ok := debug.ReadBuildInfo()
	//if ok {
	//	for _, dep := range bi.Deps {
	//		if dep.Path == "cuelang.org/go" {
	//			parts := strings.Split(dep.Version, ".")
	//			parts[len(parts)-1] = "0"
	//			return strings.Join(parts, ".")
	//		}
	//	}
	//}
	return "v0.14.1"
})

type (
	File     = cuemodfile.File
	Dep      = cuemodfile.Dep
	Language = cuemodfile.Language
)

type OldFile struct{}
