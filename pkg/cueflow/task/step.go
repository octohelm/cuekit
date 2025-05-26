package task

import (
	"iter"

	"cuelang.org/go/cue"
)

func Steps(value cue.Value) iter.Seq2[cue.Value, error] {
	return func(yield func(item cue.Value, err error) bool) {
		v := value.LookupPath(cue.ParsePath("steps"))

		list, err := v.List()
		if err != nil {
			if !yield(cue.Value{}, err) {
				return
			}
			return
		}

		for list.Next() {
			if !yield(list.Value(), nil) {
				return
			}
		}
	}
}
