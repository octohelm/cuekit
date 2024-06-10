package modfile

import cuemodfile "cuelang.org/go/mod/modfile"

const CueVersion = "v0.9.0"

type File = cuemodfile.File
type Dep = cuemodfile.Dep
type Language = cuemodfile.Language

type OldFile struct {
}
