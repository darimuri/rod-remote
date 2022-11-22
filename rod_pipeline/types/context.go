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
	m            map[string]interface{}
}

func (c *PipelineContext) PushPage(pg *rod.Page) {
	c.pageStack.Push(c.p)
	c.p = pg
}

func (c *PipelineContext) PopPage() error {
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

func (c *PipelineContext) ElementStackEmpty() bool {
	return len(c.elementStack) == 0
}

func (c *PipelineContext) PushElement(e *rod.Element) {
	c.elementStack.Push(e)
}

func (c *PipelineContext) PopElement() error {
	_, ok := c.elementStack.Pop()
	if !ok {
		return errors.New("element stack is empty")
	}

	return nil
}

func (c *PipelineContext) Set(k string, v interface{}) {
	c.m[k] = v
}

func (c *PipelineContext) Get(k string) interface{} {
	return c.m[k]
}

func NewContext(p *rod.Page) *PipelineContext {
	return &PipelineContext{
		p:            p,
		pageStack:    Stack[*rod.Page]{},
		elementStack: Stack[*rod.Element]{},
		m:            make(map[string]interface{}, 0),
	}
}
