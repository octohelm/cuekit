package modfile

import cuemodfile "cuelang.org/go/mod/modfile"

const CueVersion = "v0.9.0"

type (
	File     = cuemodfile.File
	Dep      = cuemodfile.Dep
	Language = cuemodfile.Language
)

type OldFile struct{}
