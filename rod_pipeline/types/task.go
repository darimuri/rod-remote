package types

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type ITask interface {
	Do(pc *PipelineContext) error
	Desc() string
}

type ConditionalTask struct {
	f    ConditionalFunc
	desc string
}

func (t *ConditionalTask) Do(pc *PipelineContext) (bool, error) {
	return t.f(pc)
}

func (t *ConditionalTask) Desc() string {
	return t.desc
}

func (t *ConditionalTask) DoReference() interface{} {
	return t.f
}

func NewConditionalTask(f ConditionalFunc, desc string) ConditionalTask {
	return ConditionalTask{f: f, desc: desc}
}

type OpFunc func(pc *PipelineContext) error

type ConditionalFunc func(pc *PipelineContext) (bool, error)
type EachElementFunc func(pc *PipelineContext, el *rod.Element) (bool, error)

type DialogHandlerFunc func(wait func() *proto.PageJavascriptDialogOpening, handle func(*proto.PageHandleJavaScriptDialog) error)
