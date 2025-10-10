package graph

import (
	"iter"

	"cuelang.org/go/cue"
	cuedep "cuelang.org/go/cue/dep"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/octohelm/cuekit/pkg/cuepath"
	"github.com/octohelm/cuekit/pkg/cueutil"
)

type Resolver struct {
	CreateNode func(n Node) Node
	IsPrefix   func(p cue.Path) (bool, cue.Path)

	root cue.Value

	nodes map[jsontext.Pointer]Node
}

func (r *Resolver) Init(v cue.Value) error {
	r.root = v

	if err := r.scanNodes(v, map[string]bool{}); err != nil {
		return err
	}

	for _, n := range r.nodes {
		if err := r.resolveDepsOf(n); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) Nodes() iter.Seq[Node] {
	return func(yield func(Node) bool) {
		for _, n := range r.nodes {
			if !yield(n) {
				return
			}
		}
	}
}

func (r *Resolver) addNode(n Node) {
	if r.nodes == nil {
		r.nodes = map[jsontext.Pointer]Node{}
	}

	r.nodes[cuepath.AsJSONPointer(n.Path())] = n
}

func (r *Resolver) resolveDepsOf(target Node) error {
	for path, err := range r.referencePathsOf(r.root.LookupPath(target.Path()), map[string]bool{}) {
		if err != nil {
			return err
		}

		// skip in-node
		if cuepath.Prefix(path, target.Path()) {
			continue
		}

		for _, n := range r.nodes {
			if cuepath.Prefix(path, n.Path()) {
				target.(NodeAccessor).AddDep(n)
				break
			}
		}
	}

	return nil
}

func (r *Resolver) scanNodes(v cue.Value, scanned map[string]bool) error {
	for fv, err := range cueutil.Fields(v,
		cue.Hidden(true),
		cue.Optional(true),
		cue.Definitions(false),
	) {
		if err != nil {
			return err
		}

		p := fv.Path().String()

		if _, ok := scanned[p]; !ok {
			scanned[p] = true

			if name, ok := r.isNode(fv); ok {
				r.addNode(r.createNode(&node{name: name, path: fv.Path()}))
				// don't scan deep
				continue
			}

			if err := r.scanNodes(fv, scanned); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Resolver) isNode(v cue.Value) (string, bool) {
	remain := v.Path()

	if isPrefix := r.IsPrefix; isPrefix != nil {
		ok, remainPath := isPrefix(v.Path())
		if !ok {
			return "", false
		}
		remain = remainPath
	}

	if cv := v.LookupPath(TypePath); cv.Exists() {
		if cuepath.Contains(remain, func(sel cue.Selector, i int) bool {
			return sel.Type() == cue.IndexLabel
		}) {
			return "", false
		}

		n, _ := cv.String()
		return n, true
	}
	return "", false
}

func (r *Resolver) referencePathsOf(target cue.Value, resolved map[string]bool) iter.Seq2[cue.Path, error] {
	return func(yield func(cue.Path, error) bool) {
		p := target.Path().String()

		emitted := make(map[string]bool)

		uniqueYield := func(p cue.Path, err error) bool {
			if err != nil {
				return yield(p, err)
			}

			pp := p.String()
			if emitted[pp] {
				return true
			}
			emitted[p.String()] = true

			return yield(p, err)
		}

		if _, ok := resolved[p]; !ok {
			resolved[p] = true

			for v, err := range r.values(target) {
				if err != nil {
					yield(cue.MakePath(), err)
					return
				}

				if cuepath.Contains(v.Path(), func(sel cue.Selector, i int) bool {
					return sel.IsDefinition() && i == 0
				}) {
					// skip definition
					continue
				}

				hasDep := false

				for p0, err := range cuedep.Deps(v) {
					if err != nil {
						yield(cue.MakePath(), err)
						return
					}

					if cuepath.Contains(p0, func(sel cue.Selector, i int) bool {
						return sel.IsDefinition() && i == 0
					}) {
						// skip definition
						continue
					}

					hasDep = true
					if !uniqueYield(p0, nil) {
						return
					}

					if !cuepath.Prefix(target.Path(), p0) {
						for p1, err := range r.referencePathsOf(r.root.LookupPath(p0), resolved) {
							if !uniqueYield(p1, err) {
								return
							}
						}
					}
				}

				if !hasDep {
					for p1, err := range r.referencePathsFromReferencePath(v, target, resolved) {
						if !uniqueYield(p1, err) {
							return
						}
					}
				}
			}
		}
	}
}

func (r *Resolver) referencePathsFromReferencePath(v cue.Value, target cue.Value, resolved map[string]bool) iter.Seq2[cue.Path, error] {
	safeExpr := func(v cue.Value) (op cue.Op, values []cue.Value) {
		defer func() {
			// ugly catch
			// FIXME until
			_ = recover()
		}()

		op, values = v.Expr()
		return op, values
	}

	return func(yield func(cue.Path, error) bool) {
		op, values := safeExpr(v)
		switch op {
		case cue.SelectorOp, cue.IndexOp:
			// to handle some case as dep
			// result: _pull.output.rootfs
			root, p0 := v.ReferencePath()
			if root.Exists() {
				if !yield(cuepath.Parent(p0), nil) {
					return
				}

				if !cuepath.Prefix(target.Path(), p0) {
					for p1, err := range r.referencePathsOf(r.root.LookupPath(p0), resolved) {
						if !yield(p1, err) {
							return
						}
					}
				}
			}
		case cue.InterpolationOp:
			for _, x := range values {
				for p1, err := range r.referencePathsFromReferencePath(x, target, resolved) {
					if !yield(p1, err) {
						return
					}
				}
			}
		}
	}
}

func (p *Resolver) values(target cue.Value) iter.Seq2[cue.Value, error] {
	return func(yield func(cue.Value, error) bool) {
		if !yield(target, nil) {
			return
		}

		for v, err := range cueutil.AllValues(
			target,
			cue.Optional(true),
			cue.Hidden(true),
			cue.Definitions(false),
		) {
			if !yield(v, err) {
				return
			}
		}

		op, values := target.Expr()
		switch op {
		case cue.AndOp, cue.OrOp:
			for _, v := range values {
				for vv, err := range p.values(v) {
					if !yield(vv, err) {
						return
					}
				}
			}
		}
	}
}

func (r *Resolver) createNode(n *node) Node {
	if r.CreateNode != nil {
		return r.CreateNode(n)
	}
	return n
}
