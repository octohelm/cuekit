package cueconvert

import (
	"testing"

	"github.com/octohelm/cuekit/pkg/cuepath"
	"github.com/octohelm/x/testing/bdd"
	"github.com/octohelm/x/testing/snapshot"
)

func TestFormat(t *testing.T) {
	b := bdd.FromT(t)

	simple := bdd.Must(Dump(
		map[string]any{
			"x": 1,
		},
		WithStrictValueMatcher(
			bdd.Must(cuepath.CompileGlobMatcher("x")),
		),
	))

	withStatic := bdd.Must(Dump(
		map[string]any{
			"kind": nil,
			"x":    1,
		},
		WithStaticValue(map[string]any{
			"kind": "X",
		}),
	))

	withTyped := bdd.Must(Dump(
		map[string]any{
			"x": 1,
		},
		WithType("X"),
	))

	asDecl := bdd.Must(Dump(
		map[string]any{
			"x": 1,
		},
		AsDecl("X"),
	))

	b.Then("match snapshot",
		bdd.MatchSnapshot(func(s *snapshot.Snapshot) {
			s.Add("simple.cue", simple)
			s.Add("with_static.cue", withStatic)
			s.Add("with_typed.cue", withTyped)
			s.Add("as_decl.cue", asDecl)
		}, "dump"),
	)
}
