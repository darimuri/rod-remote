package task

import (
	"github.com/go-rod/rod"
)

type ConditionalFunc func(p *rod.Page) (bool, error)

var _ ITask = (*If)(nil)

type If struct {
	op      ConditionalFunc
	iftrue  *Tasks
	iffalse *Tasks
}

func (c *If) Do(p *rod.Page) error {
	cond, err := c.op(p)
	if err != nil {
		return err
	}

	if cond {
		for _, t := range c.iftrue.tasks {
			if errT := t.Do(p); errT != nil {
				return errT
			}
		}
	} else {
		for _, t := range c.iffalse.tasks {
			if errT := t.Do(p); errT != nil {
				return errT
			}
		}
	}

	return nil
}

var _ ITask = (*While)(nil)

type While struct {
	op       ConditionalFunc
	iftrue   *Tasks
	iffalse  *Tasks
	maxRetry int
}

func (c *While) Do(p *rod.Page) error {
	for i := 0; i < c.maxRetry; i++ {
		cond, err := c.op(p)
		if err != nil {
			return err
		}

		if cond {
			for _, t := range c.iftrue.tasks {
				if errT := t.Do(p); errT != nil {
					return errT
				}
			}
			break
		} else {
			for _, t := range c.iffalse.tasks {
				if errT := t.Do(p); errT != nil {
					return errT
				}
			}
		}
	}

	return nil
}

func Has(selector string) ConditionalFunc {
	f := func(p *rod.Page) (bool, error) {
		has, _, err := p.Has(selector)
		if err != nil {
			return false, err
		}

		return has, nil
	}
	return f
}

func Visible(selector string) ConditionalFunc {
	f := func(p *rod.Page) (bool, error) {
		has, el, err := p.Has(selector)
		if err != nil {
			return false, err
		}
		if !has {
			return false, nil
		}

		return el.Visible()
	}
	return f
}
