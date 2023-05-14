package task

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/darimuri/rod-remote/rod_pipeline/types"
	"github.com/darimuri/rod-remote/rod_pipeline/util"
)

var _ types.ITask = (*IfTask)(nil)

type IfTask struct {
	op         types.ConditionalFunc
	trueTasks  *Tasks
	falseTasks *Tasks
	desc       string
}

func (c *IfTask) Do(pc *types.PipelineContext) error {
	cond, err := c.op(pc)

	if err != nil {
		return err
	}

	if cond {
		for _, t := range c.trueTasks.Tasks {
			log.Printf(">>> %s", t.Desc())
			if errT := t.Do(pc); errT != nil {
				return errT
			}
		}
	} else {
		for _, t := range c.falseTasks.Tasks {
			log.Printf(">>> %s", t.Desc())
			if errT := t.Do(pc); errT != nil {
				return errT
			}
		}
	}

	return nil
}

func (c *IfTask) Desc() string {
	return c.desc
}

var _ types.ITask = (*WhileTask)(nil)

type WhileTask struct {
	op         types.ConditionalFunc
	trueTasks  *Tasks
	falseTasks *Tasks
	maxRetry   int
	desc       string
}

func (c *WhileTask) Do(pc *types.PipelineContext) error {
	for i := 0; i < c.maxRetry; i++ {
		cond, err := c.op(pc)
		if err != nil {
			return err
		}

		if cond {
			log.Printf(">>> do truetasks for condition is %v (%d/%d)", cond, i, c.maxRetry)
			for _, t := range c.trueTasks.Tasks {
				log.Printf(">>>> %s", t.Desc())
				if err = t.Do(pc); err != nil {
					return err
				}
			}
			break
		} else {
			log.Printf(">>> do falsetasks for condition is %v (%d/%d)", cond, i, c.maxRetry)
			for _, t := range c.falseTasks.Tasks {
				log.Printf(">>>> %s", t.Desc())
				if err = t.Do(pc); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *WhileTask) Desc() string {
	return c.desc
}

var _ types.ITask = (*ForEachElementTask)(nil)

type ForEachElementTask struct {
	selector string
	tasks    *Tasks
	desc     string
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

		for _, t := range c.tasks.Tasks {
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

func (c *ForEachElementTask) Desc() string {
	return c.desc
}

func If(op types.ConditionalTask, trueTasks []types.ITask, falseTasks []types.ITask) *IfTask {
	return &IfTask{
		op:         op.Do,
		trueTasks:  NewTasks(trueTasks...),
		falseTasks: NewTasks(falseTasks...),
		desc:       fmt.Sprintf("if %s (%s)", op.Desc(), util.FunctionName(op.DoReference())),
	}
}

func While(op types.ConditionalTask, trueTasks []types.ITask, falseTasks []types.ITask, maxRetry int) *WhileTask {
	return &WhileTask{
		op:         op.Do,
		trueTasks:  NewTasks(trueTasks...),
		falseTasks: NewTasks(falseTasks...),
		maxRetry:   maxRetry,
		desc:       fmt.Sprintf("while %s (%s,%d)", op.Desc(), util.FunctionName(op.DoReference()), maxRetry),
	}
}

func WhileUntilHas(selector string, maxRetry int, delay time.Duration) *WhileTask {
	op := func(pc *types.PipelineContext) (bool, error) {
		has, _, err := pc.Query().Has(selector)
		return has, err
	}
	iftrue := NewTasks()
	iffalse := NewTasks(Sleep(delay))

	return &WhileTask{
		op:         op,
		trueTasks:  iftrue,
		falseTasks: iffalse,
		maxRetry:   maxRetry,
		desc:       fmt.Sprintf("while until has %q", selector),
	}
}

func ForEachElement(selector string, tasks []types.ITask) *ForEachElementTask {
	return &ForEachElementTask{
		selector: selector,
		tasks:    NewTasks(tasks...),
		desc:     fmt.Sprintf("foreach %q", selector),
	}
}
