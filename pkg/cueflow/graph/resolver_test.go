package graph

import (
	"fmt"
	"os"
	"testing"

	"github.com/octohelm/x/testing/bdd"
	"github.com/octohelm/x/testing/snapshot"

	"github.com/octohelm/cuekit/pkg/cueutil"
)

func FuzzResolver(f *testing.F) {
	f.Add("in_loop")
	f.Add("from_slice")
	f.Add("sub_field")
	f.Add("result_as_field_default")

	f.Fuzz(func(t *testing.T, c string) {
		b := bdd.FromT(t)

		root := bdd.Must(cueutil.BuildFile(bdd.Must(os.ReadFile(fmt.Sprintf("./testdata/%s.cue", c)))))
		fmt.Println(root)

		r := &Resolver{}
		b.Then("init success",
			bdd.NoError(r.Init(root)),
		)

		for n := range r.Nodes() {
			fmt.Println("NODE", n.Path())
			for d := range n.Deps() {
				fmt.Println(" DEP", d.Path())
			}
		}

		b.Then("resolve graph",
			bdd.MatchSnapshot(func(s *snapshot.Snapshot) {
				g := ToD2Graph(Collect(r.Nodes()))
				s.Add("g.d2", g)
				fmt.Println(bdd.Must(ToKrokiURI(g)))
			}, c),
		)
	})
}
