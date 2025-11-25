package graph

import (
	"fmt"
	"os"
	"testing"

	"cuelang.org/go/cue"

	. "github.com/octohelm/x/testing/v2"

	"github.com/octohelm/cuekit/pkg/cueutil"
)

func FuzzResolver(f *testing.F) {
	f.Add("in_loop")
	f.Add("from_slice")
	f.Add("sub_field")
	f.Add("result_as_field_default")

	f.Fuzz(func(t *testing.T, c string) {
		root := MustValue(t, func() (cue.Value, error) {
			data, err := os.ReadFile(fmt.Sprintf("./testdata/%s.cue", c))
			if err != nil {
				return cue.Value{}, err
			}
			return cueutil.BuildFile(data)
		})

		r := &Resolver{}

		t.Run(fmt.Sprintf("resolve graph from %s", c), func(t *testing.T) {
			Then(t, "success",
				ExpectMust(func() error {
					return r.Init(root)
				}),

				ExpectMustValue(func() (Snapshot, error) {
					for n := range r.Nodes() {
						fmt.Printf("NODE %s\n", n.Path())
						for d := range n.Deps() {
							fmt.Printf("  DEP %s\n", d.Path())
						}
					}

					g := ToD2Graph(Collect(r.Nodes()))

					if uri, err := ToKrokiURI(g); err == nil {
						fmt.Printf("Kroki URI: %s\n", uri)
					}

					return SnapshotOf(
						SnapshotFileFromRaw("g.d2", g),
					), nil
				}, MatchSnapshot(c)),
			)
		})
	})
}
