package runner

import (
	"cuelang.org/go/cue"
)

var (
	PathControl = cue.ParsePath("$$control.name")

	PathDep = cue.ParsePath("$dep")
	PathOk  = cue.ParsePath("$ok")
)
