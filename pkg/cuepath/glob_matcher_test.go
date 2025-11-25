package cuepath

import (
	"testing"

	"cuelang.org/go/cue"

	"github.com/octohelm/x/cmp"
	. "github.com/octohelm/x/testing/v2"
)

func TestPathMatcher(t *testing.T) {
	m := MustValue(t, func() (GlobMatcher, error) {
		return CompileGlobMatcher(`"*"."{kind,type}"`)
	})

	t.Run("Match", func(t *testing.T) {
		Then(t, "glob matcher should work as expected",
			// 匹配成功的情况
			Expect(m.Match(cue.ParsePath("x.kind")), Be(cmp.True())),
			Expect(m.Match(cue.ParsePath("a.b.c.kind")), Be(cmp.True())),
			// 匹配失败的情况
			Expect(m.Match(cue.ParsePath("a.b")), Be(cmp.False())),
		)
	})
}
