package types

import (
	"github.com/go-rod/rod"
)

type PageStack []*rod.Page

func (s *PageStack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *PageStack) Push(p *rod.Page) {
	*s = append(*s, p)
}

func (s *PageStack) Pop() (*rod.Page, bool) {
	if s.IsEmpty() {
		return nil, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}
