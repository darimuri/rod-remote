package task

import (
	"time"

	"github.com/darimuri/rod-remote/rod_pipeline/types"
)

var _ types.ITask = (*Task)(nil)

type Task struct {
	op types.OpFunc
}

func (t *Task) Do(pc *types.PipelineContext) error {
	return t.op(pc)
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

func (t *Tasks) Do(pc *types.PipelineContext) error {
	for _, task := range t.Tasks {
		if err := task.Do(pc); err != nil {
			return err
		}
	}

	return nil
}

func (t *Tasks) Open(url string) *Tasks {
	t.Append(Open(url))

	return t
}

func (t *Tasks) WaitLoad() *Tasks {
	t.Append(WaitLoad())

	return t
}

func (t *Tasks) WaitIdle(dur time.Duration) *Tasks {
	t.Append(WaitIdle(dur))

	return t
}

func (t *Tasks) Click(selector string, handler types.DialogHandlerFunc) *Tasks {
	t.Append(Click(selector, handler))

	return t
}

func (t *Tasks) Tap(selector string, handler types.DialogHandlerFunc) *Tasks {
	t.Append(Tap(selector, handler))

	return t
}

func (t *Tasks) Input(selector string, str string) *Tasks {
	t.Append(Input(selector, str))

	return t
}

func (t *Tasks) Reload() *Tasks {
	t.Append(Reload())

	return t
}

func (t *Tasks) Sleep(dur time.Duration) *Tasks {
	t.Append(Sleep(dur))

	return t
}

func (t *Tasks) Stop(message string) *Tasks {
	t.Append(Stop(message))

	return t
}

func (t *Tasks) Custom(c func(pc *types.PipelineContext) error) *Tasks {
	t.Append(Custom(c))

	return t
}

func (t *Tasks) ForEach(selector string, ef types.EachElementFunc) *Tasks {
	t.Append(ForEach(selector, ef))
	return t
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
