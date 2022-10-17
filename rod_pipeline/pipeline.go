package rodpipeline

import (
	"github.com/go-rod/rod"

	"github.com/darimuri/rod-remote/rod_pipeline/task"
	"github.com/darimuri/rod-remote/rod_pipeline/types"
)

type Pipeline struct {
	p *rod.Page
	*task.Tasks
}

func (p Pipeline) Run() error {
	return p.Tasks.Do(p.p)
}

func NewPipeline(p *rod.Page) *Pipeline {
	return &Pipeline{p: p, Tasks: &task.Tasks{}}
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
