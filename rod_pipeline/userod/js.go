package userod

import (
	"fmt"
	"strings"

	"github.com/go-rod/rod"
)

func EvalEventScript(el *rod.Element, event string) (bool, error) {
	js, errJs := el.Attribute(event)
	if errJs != nil {
		return false, errJs
	} else if js == nil {
		return false, nil
	}

	myJs := strings.ReplaceAll(*js, "return false;", "")
	_, errEval := el.Page().Eval(fmt.Sprintf("() => %s", myJs))
	return true, errEval
}
