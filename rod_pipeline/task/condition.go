package task

import (
	"github.com/go-rod/rod"

	"github.com/darimuri/rod-remote/rod_pipeline/types"
)

var _ types.ITask = (*IfTask)(nil)

type IfTask struct {
	op      types.ConditionalFunc
	iftrue  *Tasks
	iffalse *Tasks
}

func (c *IfTask) Do(p *rod.Page) error {
	cond, err := c.op(p)
	if err != nil {
		return err
	}

	if cond {
		for _, t := range c.iftrue.Tasks {
			if errT := t.Do(p); errT != nil {
				return errT
			}
		}
	} else {
		for _, t := range c.iffalse.Tasks {
			if errT := t.Do(p); errT != nil {
				return errT
			}
		}
	}

	return nil
}

var _ types.ITask = (*WhileTask)(nil)

type WhileTask struct {
	op       types.ConditionalFunc
	iftrue   *Tasks
	iffalse  *Tasks
	maxRetry int
}

func (c *WhileTask) Do(p *rod.Page) error {
	for i := 0; i < c.maxRetry; i++ {
		cond, err := c.op(p)
		if err != nil {
			return err
		}

		if cond {
			for _, t := range c.iftrue.Tasks {
				if errT := t.Do(p); errT != nil {
					return errT
				}
			}
			break
		} else {
			for _, t := range c.iffalse.Tasks {
				if errT := t.Do(p); errT != nil {
					return errT
				}
			}
		}
	}

	return nil
}

func Has(selector string) types.ConditionalFunc {
	f := func(p *rod.Page) (bool, error) {
		has, _, err := p.Has(selector)
		if err != nil {
			return false, err
		}

		return has, nil
	}
	return f
}

func Visible(selector string) types.ConditionalFunc {
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

func If(op types.ConditionalFunc, trueTasks []types.ITask, falseTasks []types.ITask) *IfTask {
	iftrue := NewTasks(trueTasks...)
	iffalse := NewTasks(falseTasks...)

	return &IfTask{op: op, iftrue: iftrue, iffalse: iffalse}
}

func While(op types.ConditionalFunc, trueTasks []types.ITask, falseTasks []types.ITask, maxRetry int) *WhileTask {
	iftrue := NewTasks(trueTasks...)
	iffalse := NewTasks(falseTasks...)

	conditional := &WhileTask{op: op, iftrue: iftrue, iffalse: iffalse, maxRetry: maxRetry}
	return conditional
}
