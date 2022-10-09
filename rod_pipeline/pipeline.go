package rodpipeline

import (
	"github.com/go-rod/rod"

	"github.com/darimuri/rod-remote/rod_pipeline/task"
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
