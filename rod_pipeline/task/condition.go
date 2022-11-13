package task

import (
	"strings"

	"github.com/darimuri/rod-remote/rod_pipeline/types"
)

var _ types.ITask = (*IfTask)(nil)

type IfTask struct {
	op      types.ConditionalFunc
	iftrue  *Tasks
	iffalse *Tasks
}

func (c *IfTask) Do(pc *types.PipelineContext) error {
	cond, err := c.op(pc)
	if err != nil {
		return err
	}

	if cond {
		for _, t := range c.iftrue.Tasks {
			if errT := t.Do(pc); errT != nil {
				return errT
			}
		}
	} else {
		for _, t := range c.iffalse.Tasks {
			if errT := t.Do(pc); errT != nil {
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

func (c *WhileTask) Do(pc *types.PipelineContext) error {
	for i := 0; i < c.maxRetry; i++ {
		cond, err := c.op(pc)
		if err != nil {
			return err
		}

		if cond {
			for _, t := range c.iftrue.Tasks {
				if errT := t.Do(pc); errT != nil {
					return errT
				}
			}
			break
		} else {
			for _, t := range c.iffalse.Tasks {
				if errT := t.Do(pc); errT != nil {
					return errT
				}
			}
		}
	}

	return nil
}

func Has(selector string) types.ConditionalFunc {
	f := func(pc *types.PipelineContext) (bool, error) {
		has, _, err := pc.Query().Has(selector)
		if err != nil {
			return false, err
		}

		return has, nil
	}
	return f
}

func ContainsText(selector, text string) types.ConditionalFunc {
	f := func(pc *types.PipelineContext) (bool, error) {
		has, el, err := pc.Query().Has(selector)
		if err != nil {
			return false, err
		}
		if false == has {
			return has, nil
		}

		s, errText := el.Text()
		if errText != nil {
			return false, errText
		}

		return strings.Contains(strings.TrimSpace(s), text), nil
	}

	return f
}

func Visible(selector string) types.ConditionalFunc {
	f := func(pc *types.PipelineContext) (bool, error) {
		has, el, err := pc.Query().Has(selector)
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
