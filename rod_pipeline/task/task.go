package task

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"

	"github.com/darimuri/rod-remote/rod_pipeline/types"
)

var _ types.ITask = (*Task)(nil)

type Task struct {
	op types.OpFunc
}

func (t *Task) Do(c *types.PipelineContext) error {
	return t.op(c)
}

type Tasks struct {
	Tasks []types.ITask
}

func NewTasks(task ...types.ITask) *Tasks {
	t := &Tasks{}
	t.Append(task...)
	return t
}

func (t *Tasks) Append(task ...types.ITask) {
	t.Tasks = append(t.Tasks, task...)
}

func (t *Tasks) Do(c *types.PipelineContext) error {
	for _, task := range t.Tasks {
		if err := task.Do(c); err != nil {
			return err
		}
	}

	return nil
}

func (t *Tasks) Open(url string) *Tasks {
	t.Append(Open(url))

	return t
}

func Open(url string) *Task {
	f := func(ctx *types.PipelineContext) error {
		return ctx.P.Navigate(url)
	}
	task := &Task{op: f}
	return task
}

func (t *Tasks) WaitLoad() *Tasks {
	t.Append(WaitLoad())

	return t
}

func WaitLoad() *Task {
	f := func(ctx *types.PipelineContext) error {
		return ctx.P.WaitLoad()
	}
	task := &Task{op: f}
	return task
}

func (t *Tasks) WaitIdle(dur time.Duration) *Tasks {
	t.Append(WaitIdle(dur))

	return t
}

func WaitIdle(dur time.Duration) *Task {
	f := func(ctx *types.PipelineContext) error {
		return ctx.P.WaitIdle(dur)
	}
	task := &Task{op: f}
	return task
}

func (t *Tasks) Click(selector string, handler types.DialogHandlerFunc) *Tasks {
	t.Append(Click(selector, handler))

	return t
}

func Click(selector string, handler types.DialogHandlerFunc) *Task {
	f := func(ctx *types.PipelineContext) error {
		el, err := ctx.P.Element(selector)
		if err != nil {
			return err
		}

		if err = el.Hover(); err != nil {
			return err
		}

		if handler != nil {
			timeout := ctx.P.Timeout(time.Minute)
			wait, handle := timeout.HandleDialog()
			go handler(wait, handle)
		}

		return el.Click(proto.InputMouseButtonLeft, 1)
	}
	task := &Task{op: f}
	return task
}

func (t *Tasks) Tap(selector string, handler types.DialogHandlerFunc) *Tasks {
	t.Append(Tap(selector, handler))

	return t
}

func Tap(selector string, handler types.DialogHandlerFunc) *Task {
	f := func(ctx *types.PipelineContext) error {
		el, err := ctx.P.Element(selector)
		if err != nil {
			return err
		}

		if err = el.Hover(); err != nil {
			return err
		}

		if handler != nil {
			timeout := ctx.P.Timeout(time.Minute)
			wait, handle := timeout.HandleDialog()
			go handler(wait, handle)
		}

		return el.Tap()
	}
	task := &Task{op: f}
	return task
}

func (t *Tasks) Custom(c func(p *rod.Page) error) *Tasks {
	t.Append(Custom(c))

	return t
}

func Custom(c func(p *rod.Page) error) *Task {
	f := func(ctx *types.PipelineContext) error {
		return c(ctx.P)
	}
	task := &Task{op: f}
	return task
}

func RemoveClass(selector string, class string) *Task {
	f := func(ctx *types.PipelineContext) error {
		el, err := ctx.P.Element(selector)
		if err != nil {
			return err
		}

		_, errEval := el.Eval(fmt.Sprintf(`() => this.classList.remove('%s')`, class))
		if errEval != nil {
			return errEval
		}

		return nil
	}
	task := &Task{op: f}
	return task
}

func AddClass(selector string, class string) *Task {
	f := func(ctx *types.PipelineContext) error {
		el, err := ctx.P.Element(selector)
		if err != nil {
			return err
		}

		_, errEval := el.Eval(fmt.Sprintf(`() => this.classList.add('%s')`, class))
		if errEval != nil {
			return errEval
		}

		return nil
	}
	task := &Task{op: f}
	return task
}

func (t *Tasks) Input(selector string, str string) *Tasks {
	t.Append(Input(selector, str))

	return t
}

func Input(selector string, str string) *Task {
	f := func(ctx *types.PipelineContext) error {
		el, err := ctx.P.Element(selector)
		if err != nil {
			return err
		}

		if err = el.SelectAllText(); err != nil {
			return err
		}

		if err = el.Input(str); err != nil {
			return err
		}

		return nil
	}
	task := &Task{op: f}
	return task
}

func Type(keys ...input.Key) *Task {
	f := func(ctx *types.PipelineContext) error {
		if err := ctx.P.Keyboard.Type(keys...); err != nil {
			return err
		}

		return nil
	}
	task := &Task{op: f}
	return task

}

func (t *Tasks) Reload() *Tasks {
	t.Append(Reload())

	return t
}

func Reload() *Task {
	f := func(ctx *types.PipelineContext) error {
		return ctx.P.Reload()
	}
	task := &Task{op: f}
	return task
}

func (t *Tasks) Sleep(dur time.Duration) *Tasks {
	t.Append(Sleep(dur))

	return t
}

func Sleep(dur time.Duration) *Task {
	f := func(ctx *types.PipelineContext) error {
		time.Sleep(dur)
		return nil
	}
	task := &Task{op: f}
	return task
}

func (t *Tasks) Stop(message string) *Tasks {
	t.Append(Stop(message))

	return t
}

func Stop(message string) *Task {
	f := func(ctx *types.PipelineContext) error {
		return errors.New(message)
	}
	task := &Task{op: f}
	return task
}

func (t *Tasks) ForEachElement(selector string, ef types.EachElementFunc) *Tasks {
	t.Append(ForEach(selector, ef))
	return t
}

func ForEach(selector string, ef types.EachElementFunc) *Task {
	f := func(ctx *types.PipelineContext) error {
		els, err := ctx.P.Elements(selector)
		if err != nil {
			return err
		}
		for _, el := range els {
			if stop, errEl := ef(el); errEl != nil {
				return errEl
			} else if stop {
				break
			}
		}
		return nil
	}
	return &Task{op: f}
}

func (t *Tasks) If(op types.ConditionalFunc, trueTasks, falseTasks []types.ITask) *Tasks {
	t.Append(If(op, trueTasks, falseTasks))
	return t
}

func (t *Tasks) While(op types.ConditionalFunc, trueTasks, falseTasks []types.ITask, maxRetry int) *Tasks {
	conditional := While(op, trueTasks, falseTasks, maxRetry)

	t.Append(conditional)
	return t
}
