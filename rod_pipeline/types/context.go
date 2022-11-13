package types

import (
	"errors"

	"github.com/go-rod/rod"

	"github.com/darimuri/rod-remote/rod_pipeline/userod"
)

type PipelineContext struct {
	p            *rod.Page
	pageStack    Stack[*rod.Page]
	elementStack Stack[*rod.Element]
}

func (c *PipelineContext) Push(pg *rod.Page) {
	c.pageStack.Push(c.p)
	c.p = pg
}

func (c *PipelineContext) Pop() error {
	pp, ok := c.pageStack.Pop()
	if ok {
		c.p = pp
	} else {
		return errors.New("page stack is empty")
	}

	return nil
}

func (c *PipelineContext) Page() *rod.Page {
	return c.p
}

func (c *PipelineContext) Query() userod.Query {
	if c.elementStack.IsEmpty() {
		return userod.NewQuery(c.p, nil)
	}
	return userod.NewQuery(c.p, c.elementStack.Last())
}

func NewContext(p *rod.Page) *PipelineContext {
	return &PipelineContext{
		p:            p,
		pageStack:    Stack[*rod.Page]{},
		elementStack: Stack[*rod.Element]{},
	}
}
