package types

type Stack[T interface{}] []T

func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack[T]) Push(p T) {
	*s = append(*s, p)
}

func (s *Stack[T]) Pop() (T, bool) {
	var e T
	if s.IsEmpty() {
		return e, false
	} else {
		index := len(*s) - 1 // Get the index of the top most element.
		e := (*s)[index]     // Index into the slice and obtain the element.
		*s = (*s)[:index]    // Remove it from the stack by slicing it off.
		return e, true
	}
}

func (s *Stack[T]) Last() T {
	return (*s)[len(*s)-1]
}
