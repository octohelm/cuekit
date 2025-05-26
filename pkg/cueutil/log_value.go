package cueutil

import (
	"fmt"
	"log/slog"

	"cuelang.org/go/cue"
	cueformat "cuelang.org/go/cue/format"
	encodingcue "github.com/octohelm/cuekit/pkg/encoding/cue"
)

func AsLogValue(v any) slog.LogValuer {
	return &logValue{v: v}
}

type logValue struct {
	v any
}

func (c *logValue) LogValue() slog.Value {
	switch x := c.v.(type) {
	case cue.Value:
		node := x.Syntax(
			cue.Concrete(false), // allow incomplete values
			cue.DisallowCycles(true),
			cue.Docs(true),
			cue.All(),
		)

		data, _ := cueformat.Node(node, cueformat.Simplify())
		return slog.StringValue(string(data))
	default:
		data, err := encodingcue.Marshal(c.v)
		if err != nil {
			panic(fmt.Errorf("encoding failed: %s, %v", err, c.v))
		}
		return slog.AnyValue(data)
	}
}
