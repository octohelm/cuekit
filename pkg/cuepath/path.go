package cuepath

import (
	"slices"

	"cuelang.org/go/cue"

	slicesx "github.com/octohelm/x/slices"
)

func Contains(path cue.Path, match func(sel cue.Selector, i int) bool) bool {
	n := len(path.Selectors())
	if n == 0 {
		return false
	}

	for i, sel := range path.Selectors() {
		if match(sel, i) {
			return true
		}
	}

	return false
}

func TrimPrefix(target cue.Path, prefix cue.Path) cue.Path {
	selectors := target.Selectors()
	prefixSelectors := prefix.Selectors()

	if len(selectors) < len(prefixSelectors) {
		return target
	}

	until := -1

	for i, x := range prefixSelectors {
		if x.String() != selectors[i].String() {
			break
		}
		until = i
	}

	return cue.MakePath(selectors[until+1:]...)
}

func Prefix(target cue.Path, prefix cue.Path) bool {
	selectors := target.Selectors()
	prefixSelectors := prefix.Selectors()

	if len(selectors) < len(prefixSelectors) {
		return false
	}

	for i, x := range prefixSelectors {
		if x.String() != selectors[i].String() {
			return false
		}
	}

	return true
}

func Suffix(target cue.Path, suffix cue.Path) bool {
	selectors := target.Selectors()
	suffixSelectors := suffix.Selectors()

	if len(selectors) < len(suffixSelectors) {
		return false
	}

	for i := 1; i <= len(suffixSelectors); i++ {
		if selectors[len(selectors)-i].String() != suffixSelectors[len(suffixSelectors)-i].String() {
			return false
		}
	}

	return true
}

func Same(target cue.Path, another cue.Path) bool {
	selectors := target.Selectors()
	anotherSelectors := another.Selectors()

	if len(selectors) != len(anotherSelectors) {
		return false
	}

	for i, x := range anotherSelectors {
		if x.String() != selectors[i].String() {
			return false
		}
	}

	return true
}

func Parent(path cue.Path) cue.Path {
	selectors := path.Selectors()
	if len(selectors) == 0 {
		return path
	}
	return cue.MakePath(selectors[0 : len(selectors)-1]...)
}

func Join(paths ...cue.Path) cue.Path {
	return cue.MakePath(
		slices.Concat(
			slicesx.Map(paths, func(p cue.Path) []cue.Selector {
				return p.Selectors()
			})...,
		)...)
}
