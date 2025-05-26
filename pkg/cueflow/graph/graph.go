package graph

import (
	"bytes"
	"cmp"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"iter"
	"slices"

	"github.com/octohelm/cuekit/pkg/cuepath"
)

type Graph interface {
	Nodes() iter.Seq[Node]
}

type Elem interface {
	String() string
}

type CanShape interface {
	Shape() string
}

func Collect(nodes iter.Seq[Node]) iter.Seq[Elem] {
	return func(yield func(Elem) bool) {
		for _, task := range slices.SortedFunc(nodes, func(a Node, b Node) int {
			return cmp.Compare(cuepath.AsJSONPointer(a.Path()), cuepath.AsJSONPointer(b.Path()))
		}) {
			path := cuepath.AsJSONPointer(task.Path())

			n := &d2node{At: string(path) + "\n" + task.Name()}

			if canShape, ok := task.(CanShape); ok {
				n.Shape = canShape.Shape()
			}

			if !yield(n) {
				return
			}

			for _, dep := range slices.Sorted(
				func(yield func(jp string) bool) {
					for dep := range task.Deps() {
						if !yield(string(cuepath.AsJSONPointer(dep.Path())) + "\n" + dep.Name()) {
							return
						}
					}
				},
			) {
				if !yield(&d2link{From: dep, To: n.At}) {
					return
				}
			}
		}
	}
}

type d2node struct {
	At    string
	Shape string
}

func (d *d2node) String() string {
	def := ""
	if d.Shape != "" {
		def = fmt.Sprintf(`{ shape: %s }`, d.Shape)
	}

	return fmt.Sprintf(`
%q%s`, d.At, def)
}

type d2link struct {
	From string
	To   string
}

func (d *d2link) String() string {
	return fmt.Sprintf(`
%q -> %q`, d.From, d.To)
}

func ToD2Graph(elems iter.Seq[Elem]) []byte {
	buf := bytes.NewBuffer(nil)

	_, _ = fmt.Fprintf(buf, `direction: right
`)

	for d := range elems {
		buf.WriteString(d.String())
	}

	return buf.Bytes()
}

func ToKrokiURI(g []byte) (string, error) {
	b := bytes.NewBuffer(nil)

	w, err := zlib.NewWriterLevel(b, 9)
	if err != nil {
		return "", err
	}
	if _, err := w.Write(g); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	return fmt.Sprintf("https://kroki.io/d2/svg/%s?theme=101", base64.URLEncoding.EncodeToString(b.Bytes())), nil
}
