package cueutil

import (
	"iter"
	"strings"

	"cuelang.org/go/cue"

	"github.com/octohelm/cuekit/pkg/cuepath"
)

func Items(v cue.Value) iter.Seq2[cue.Value, error] {
	return Act(func(yield func(v cue.Value) bool) error {
		cueIter, err := v.List()
		if err != nil {
			return err
		}

		for cueIter.Next() {
			if !yield(cueIter.Value()) {
				return nil
			}
		}

		return nil
	})
}

func Fields(v cue.Value, options ...cue.Option) iter.Seq2[cue.Value, error] {
	return Act(func(yield func(v cue.Value) bool) error {
		cueIter, err := v.Fields(options...)
		if err != nil {
			return err
		}

		for cueIter.Next() {
			if !yield(cueIter.Value()) {
				return nil
			}
		}

		return nil
	})
}

func AllValues(v cue.Value, options ...cue.Option) iter.Seq2[cue.Value, error] {
	return Act(func(yield func(v cue.Value) bool) error {
		switch v.Kind() {
		case cue.ListKind:
			for item, err := range Items(v) {
				if err != nil {
					return err
				}

				if !yield(item) {
					return nil
				}

				for sub, err := range AllValues(item, options...) {
					if err != nil {
						return err
					}

					if !yield(sub) {
						return nil
					}
				}
			}
		case cue.StructKind:
			for field, err := range Fields(v, options...) {
				if err != nil {
					return err
				}

				if cuepath.Contains(field.Path(), func(sel cue.Selector, i int) bool {
					return sel.Type() == cue.StringLabel && strings.HasPrefix(sel.Unquoted(), "$$")
				}) {
					continue
				}

				if !yield(field) {
					return nil
				}

				for sub, err := range AllValues(field, options...) {
					if err != nil {
						return err
					}
					if !yield(sub) {
						return nil
					}
				}
			}
		}

		return nil
	})
}
