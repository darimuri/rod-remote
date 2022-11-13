package types

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type ITask interface {
	Do(pc *PipelineContext) error
}

type OpFunc func(pc *PipelineContext) error

type ConditionalFunc func(pc *PipelineContext) (bool, error)
type EachElementFunc func(el *rod.Element) (bool, error)

type DialogHandlerFunc func(wait func() *proto.PageJavascriptDialogOpening, handle func(*proto.PageHandleJavaScriptDialog) error)
