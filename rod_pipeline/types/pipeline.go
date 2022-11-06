package types

import (
	"errors"

	"github.com/go-rod/rod"
)

type PipelineContext struct {
	P         *rod.Page
	pageStack PageStack
}

func (c *PipelineContext) Push(pg *rod.Page) {
	c.pageStack.Push(c.P)
	c.P = pg
}

func (c *PipelineContext) Pop() error {
	pp, ok := c.pageStack.Pop()
	if ok {
		c.P = pp
	} else {
		return errors.New("page stack is empty")
	}

	return nil
}

func NewContext(p *rod.Page) *PipelineContext {
	return &PipelineContext{
		P:         p,
		pageStack: PageStack{},
	}
}
