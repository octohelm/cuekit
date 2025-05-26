package graph

import (
	"iter"

	"cuelang.org/go/cue"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/octohelm/cuekit/pkg/cuepath"
)

type Node interface {
	Name() string
	Path() cue.Path
	Value() cue.Value

	Deps() iter.Seq[Node]
}

type NodeAccessor interface {
	AddDep(n Node)
}

var TypePath = cue.ParsePath("$$type.name")

type node struct {
	name  string
	value cue.Value
	deps  map[jsontext.Pointer]Node
}

var _ Node = &node{}

func (n *node) Name() string {
	return n.name
}

func (n *node) Path() cue.Path {
	return n.value.Path()
}

func (n *node) Value() cue.Value {
	return n.value
}

func (n *node) Deps() iter.Seq[Node] {
	return func(yield func(Node) bool) {
		for _, d := range n.deps {
			if !yield(d) {
				return
			}
		}
	}
}

var _ NodeAccessor = &node{}

func (n *node) AddDep(dep Node) {
	if n.deps == nil {
		n.deps = make(map[jsontext.Pointer]Node)
	}
	n.deps[cuepath.AsJSONPointer(dep.Path())] = dep
}
