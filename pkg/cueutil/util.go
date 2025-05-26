package cueutil

import "iter"

func Act[V any](do func(yield func(V) bool) error) iter.Seq2[V, error] {
	return func(yield func(V, error) bool) {
		yieldValue := func(v V) bool { return yield(v, nil) }

		if err := do(yieldValue); err != nil {
			yield(*(new(V)), err)
			return
		}
	}
}
