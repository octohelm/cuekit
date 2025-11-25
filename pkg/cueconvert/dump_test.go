package cueconvert

import (
	"testing"

	. "github.com/octohelm/x/testing/v2"

	"github.com/octohelm/cuekit/pkg/cuepath"
)

func TestFormat(t *testing.T) {
	t.Run("Dump", func(t *testing.T) {
		Then(t, "match snapshot",
			ExpectMustValue(func() (Snapshot, error) {
				m, err := cuepath.CompileGlobMatcher("x")
				if err != nil {
					return nil, err
				}

				simple, err := Dump(
					map[string]any{"x": 1},
					WithStrictValueMatcher(m),
				)
				if err != nil {
					return nil, err
				}

				withStatic, err := Dump(
					map[string]any{"kind": nil, "x": 1},
					WithStaticValue(map[string]any{"kind": "X"}),
				)
				if err != nil {
					return nil, err
				}

				withTyped, err := Dump(
					map[string]any{"x": 1},
					WithType("X"),
				)
				if err != nil {
					return nil, err
				}

				asDecl, err := Dump(
					map[string]any{"x": 1},
					AsDecl("X"),
				)
				if err != nil {
					return nil, err
				}

				return SnapshotOf(
					SnapshotFileFromRaw("simple.cue", simple),
					SnapshotFileFromRaw("with_static.cue", withStatic),
					SnapshotFileFromRaw("with_typed.cue", withTyped),
					SnapshotFileFromRaw("as_decl.cue", asDecl),
				), nil
			}, MatchSnapshot("dump")),
		)
	})
}
