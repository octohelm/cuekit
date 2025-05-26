package cueflow

import (
	"fmt"
	"sync"

	"cuelang.org/go/cue"
)

type Scope interface {
	Value() cue.Value
	LookupPath(path cue.Path) cue.Value
	FillPath(path cue.Path, value any) error
}

type scope struct {
	root cue.Value
	mu   sync.RWMutex
}

var _ Scope = &scope{}

func (c *scope) Value() cue.Value {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.root
}

func (c *scope) LookupPath(path cue.Path) cue.Value {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.root.LookupPath(path)
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
