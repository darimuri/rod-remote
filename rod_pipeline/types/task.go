package types

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type ITask interface {
	Do(c *PipelineContext) error
}

type OpFunc func(ctx *PipelineContext) error

type ConditionalFunc func(ctx *PipelineContext) (bool, error)
type EachElementFunc func(el *rod.Element) (bool, error)

type DialogHandlerFunc func(wait func() *proto.PageJavascriptDialogOpening, handle func(*proto.PageHandleJavaScriptDialog) error)
