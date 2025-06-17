package cueflow

import (
	"fmt"
	"sync"

	"cuelang.org/go/cue"
)

type Scope interface {
	Value() cue.Value

	LookupPath(path cue.Path) cue.Value
	DecodePath(path cue.Path, target any) error
	DecodePathWith(path cue.Path, with func(x cue.Value) error) error

	FillPath(path cue.Path, value any) error
}

type scope struct {
	mu   sync.RWMutex
	root cue.Value
}

var _ Scope = &scope{}

func (c *scope) Value() cue.Value {
	return c.root
}

func (c *scope) DecodePath(path cue.Path, target any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.root.LookupPath(path).Decode(target)
}

func (c *scope) LookupPath(path cue.Path) cue.Value {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.root.LookupPath(path)
}

func (c *scope) DecodePathWith(path cue.Path, fn func(v cue.Value) error) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return fn(c.root.LookupPath(path))
}

func (c *scope) FillPath(path cue.Path, v any) error {
	if _, ok := v.(cue.Value); ok {
		return fmt.Errorf("invalid value for filling %s", path)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.root = c.root.FillPath(path, v)
	return nil
}
