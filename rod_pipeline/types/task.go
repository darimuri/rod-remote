package types

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type ITask interface {
	Do(p *rod.Page) error
}

type OpFunc func(p *rod.Page) error

type ConditionalFunc func(p *rod.Page) (bool, error)

type DialogHandlerFunc func(wait func() *proto.PageJavascriptDialogOpening, handle func(*proto.PageHandleJavaScriptDialog) error)
