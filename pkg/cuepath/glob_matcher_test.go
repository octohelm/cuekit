package cuepath

import (
	"testing"

	"cuelang.org/go/cue"

	testingx "github.com/octohelm/x/testing"
	"github.com/octohelm/x/testing/bdd"
)

func TestPathMatcher(t *testing.T) {
	m := bdd.Must(CompileGlobMatcher(`"*"."{kind,type}"`))

	testingx.Expect(t, m.Match(cue.ParsePath("x.kind")), testingx.BeTrue())
	testingx.Expect(t, m.Match(cue.ParsePath("a.b.c.kind")), testingx.BeTrue())
	testingx.Expect(t, m.Match(cue.ParsePath("a.b")), testingx.BeFalse())
}
