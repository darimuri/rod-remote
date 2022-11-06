package rodpipeline

import (
	"github.com/go-rod/rod"

	"github.com/darimuri/rod-remote/rod_pipeline/task"
	"github.com/darimuri/rod-remote/rod_pipeline/types"
)

type Pipeline struct {
	c *types.PipelineContext
	*task.Tasks
}

func (p Pipeline) Run() error {
	return p.Tasks.Do(p.c)
}

func (p Pipeline) PushPage(pg *rod.Page) {
	p.c.Push(pg)
}

func (p Pipeline) PopPage() error {
	return p.c.Pop()
}

func NewPipeline(p *rod.Page) *Pipeline {
	return &Pipeline{c: types.NewContext(p), Tasks: &task.Tasks{}}
}

func Tasks(t ...types.ITask) []types.ITask {
	tasks := make([]types.ITask, 0)
	tasks = append(tasks, t...)
	return tasks
}

func Then(t ...types.ITask) []types.ITask {
	return Tasks(t...)
}

func Else(t ...types.ITask) []types.ITask {
	return Tasks(t...)
}
