package task

import (
	"errors"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/cdp"
	"github.com/go-rod/rod/lib/proto"
)

var _ ITask = (*Task)(nil)

type ITask interface {
	Do(p *rod.Page) error
}

type OpFunc func(p *rod.Page) error

type Task struct {
	op OpFunc
}

func (t *Task) Do(p *rod.Page) error {
	return t.op(p)
}

type Tasks struct {
	tasks []ITask
}

func NewTasks(task ...ITask) *Tasks {
	t := &Tasks{}
	t.Append(task...)
	return t
}

func (t *Tasks) Append(task ...ITask) {
	t.tasks = append(t.tasks, task...)
}

func (t *Tasks) Do(p *rod.Page) error {
	for _, task := range t.tasks {
		if err := task.Do(p); err != nil {
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
	f := func(p *rod.Page) error {
		return p.Navigate(url)
	}
	task := &Task{op: f}
	return task
}

func (t *Tasks) WaitLoad() *Tasks {
	t.Append(WaitLoad())

	return t
}

func WaitLoad() *Task {
	f := func(p *rod.Page) error {
		return p.WaitLoad()
	}
	task := &Task{op: f}
	return task
}

func (t *Tasks) WaitIdle(dur time.Duration) *Tasks {
	t.Append(WaitIdle(dur))

	return t
}

func WaitIdle(dur time.Duration) *Task {
	f := func(p *rod.Page) error {
		return p.WaitIdle(dur)
	}
	task := &Task{op: f}
	return task
}

func (t *Tasks) Click(selector string) *Tasks {
	t.Append(Click(selector))

	return t
}

func Click(selector string) *Task {
	f := func(p *rod.Page) error {
		el, err := p.Element(selector)
		if err != nil {
			return err
		}

		wait, handle := p.HandleDialog()
		go func() {
			wait()
			errHandle := handle(&proto.PageHandleJavaScriptDialog{Accept: false, PromptText: ""})
			if errHandle != nil {
				if cdpError, ok := errHandle.(*cdp.Error); ok {
					if cdpError.Code != -32602 { //No dialog is showing
						panic(cdpError)
					}
				} else if errHandle.Error() != "context canceled" {
					panic(errHandle)
				}
			}
		}()

		return el.Click(proto.InputMouseButtonLeft, 1)
	}
	task := &Task{op: f}
	return task
}

func (t *Tasks) Input(selector string, str string) *Tasks {
	t.Append(Input(selector, str))

	return t
}

func Input(selector string, str string) *Task {
	f := func(p *rod.Page) error {
		el, err := p.Element(selector)
		if err != nil {
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

func (t *Tasks) Reload() *Tasks {
	t.Append(Reload())

	return t
}

func Reload() *Task {
	f := func(p *rod.Page) error {
		return p.Reload()
	}
	task := &Task{op: f}
	return task
}

func (t *Tasks) Sleep(dur time.Duration) *Tasks {
	t.Append(Sleep(dur))

	return t
}

func Sleep(dur time.Duration) *Task {
	f := func(p *rod.Page) error {
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
	f := func(p *rod.Page) error {
		return errors.New(message)
	}
	task := &Task{op: f}
	return task
}

func (t *Tasks) If(op ConditionalFunc, trueTasks, falseTasks []ITask) *Tasks {
	iftrue := NewTasks(trueTasks...)
	iffalse := NewTasks(falseTasks...)

	conditional := &If{op: op, iftrue: iftrue, iffalse: iffalse}

	t.Append(conditional)
	return t
}

func (t *Tasks) While(op ConditionalFunc, trueTasks, falseTasks []ITask, maxRetry int) *Tasks {
	iftrue := NewTasks(trueTasks...)
	iffalse := NewTasks(falseTasks...)

	conditional := &While{op: op, iftrue: iftrue, iffalse: iffalse, maxRetry: maxRetry}

	t.Append(conditional)
	return t
}

func Then(task ...ITask) []ITask {
	tasks := make([]ITask, 0)
	tasks = append(tasks, task...)
	return tasks
}

func Else(task ...ITask) []ITask {
	tasks := make([]ITask, 0)
	tasks = append(tasks, task...)
	return tasks
}
