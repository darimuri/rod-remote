package task

import (
	"errors"
	"time"

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
				if err = t.Do(pc); err != nil {
					return err
				}
			}
			break
		} else {
			for _, t := range c.iffalse.Tasks {
				if err = t.Do(pc); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

var _ types.ITask = (*ForEachElementTask)(nil)

type ForEachElementTask struct {
	selector string
	each     *Tasks
}

func (c *ForEachElementTask) Do(pc *types.PipelineContext) error {
	if !pc.ElementStackEmpty() {
		return errors.New("element stack is not empty")
	}

	elements, err := pc.Query().Elements(c.selector)
	if err != nil {
		return err
	}

	for _, e := range elements {
		pc.PushElement(e)

		for _, t := range c.each.Tasks {
			if err = t.Do(pc); err != nil {
				return err
			}
		}

		if err = pc.PopElement(); err != nil {
			return err
		}
	}

	return nil
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

func WaitUntilHas(selector string, maxRetry int, delay time.Duration) *WhileTask {
	op := func(pc *types.PipelineContext) (bool, error) {
		has, _, err := pc.Query().Has(selector)
		return has, err
	}
	iftrue := NewTasks()
	iffalse := NewTasks(Sleep(delay))

	return &WhileTask{op: op, iftrue: iftrue, iffalse: iffalse, maxRetry: maxRetry}
}

func ForEachElement(selector string, tasks []types.ITask) *ForEachElementTask {
	each := NewTasks(tasks...)
	return &ForEachElementTask{selector: selector, each: each}
}
