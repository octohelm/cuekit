package runner

import (
	"fmt"
	"iter"
	"maps"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	"github.com/octohelm/cuekit/pkg/cueconvert"
	"github.com/octohelm/cuekit/pkg/cueflow"
)

type Registry interface {
	Register(t any)
	New(task cueflow.Task) (TaskRunner, error)
	Tasks() iter.Seq[Task]
}

func NewRegistry(domain string) Registry {
	return &registry{
		domain: domain,
		types:  map[string]*named{},
	}
}

type registry struct {
	domain string
	types  map[string]*named
}

func (r *registry) Register(t any) {
	tpe := reflect.TypeOf(t)
	for tpe.Kind() == reflect.Ptr {
		tpe = tpe.Elem()
	}
	r.register(tpe)
}

func (r *registry) Tasks() iter.Seq[Task] {
	return func(yield func(Task) bool) {
		for _, t := range slices.Sorted(maps.Keys(r.types)) {
			if !yield(r.types[t]) {
				return
			}
		}
	}
}

func (r *registry) New(task cueflow.Task) (TaskRunner, error) {
	if found, ok := r.types[task.Name()]; ok {
		return found.NewRunner(task)
	}
	return nil, fmt.Errorf("unknown named `%s`", task)
}

func (f *registry) register(tpe reflect.Type) {
	block := cueconvert.FromType(tpe,
		cueconvert.WithPkgPathReplaceFunc(func(pkgPath string) string {
			if f.domain == "" {
				return ""
			}
			return fmt.Sprintf("%s/%s", f.domain, filepath.Base(pkgPath))
		}),
		cueconvert.WithRegister(f.register),
	)

	pt := &named{
		tpe:          tpe,
		decl:         block,
		outputFields: map[string][]int{},
	}

	if tpe.Kind() == reflect.Struct && !strings.HasSuffix(tpe.Name(), "Interface") {
		pt.flowStruct = true
	}

	x := reflect.New(tpe).Interface()

	if _, ok := x.(cueflow.FlowTask); ok {
		// task always struct
		pt.flowStruct = true
		pt.flowTask = true

		if _, ok := x.(cueflow.BeforeAll); ok {
			pt.flowTaskBeforeAll = true
		}
	}

	if _, ok := x.(cueflow.FlowControl); ok {
		pt.flowControl = true
	}

	for _, info := range pt.decl.Fields {
		if info.AsOutput {
			pt.outputFields[info.Name] = info.Loc
		}
	}

	f.types[pt.FullName()] = pt
}
