package task

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"

	"github.com/darimuri/rod-remote/rod_pipeline/types"
	"github.com/darimuri/rod-remote/rod_pipeline/util"
)

func Open(url string) *Task {
	f := func(pc *types.PipelineContext) error {
		return pc.Page().Navigate(url)
	}
	task := &Task{op: f, desc: fmt.Sprintf("open %s", url)}
	return task
}

func WaitRequestIdle(timeout time.Duration) *Task {
	f := func(pc *types.PipelineContext) error {
		wait := pc.Page().Timeout(timeout).MustWaitRequestIdle()
		wait()
		return nil
	}
	task := &Task{op: f, desc: fmt.Sprintf("wait request idle until %s", timeout.String())}
	return task
}

func Click(selector string, handler types.DialogHandlerFunc) *Task {
	f := func(pc *types.PipelineContext) error {
		el, err := pc.Query().Element(selector)
		if err != nil {
			return err
		}

		if err = el.Hover(); err != nil {
			return err
		}

		if handler != nil {
			timeout := pc.Page().Timeout(time.Minute)
			wait, handle := timeout.HandleDialog()
			go handler(wait, handle)
		}

		return el.Click(proto.InputMouseButtonLeft, 1)
	}
	task := &Task{op: f, desc: fmt.Sprintf("click %s", selector)}
	return task
}

func Tap(selector string, handler types.DialogHandlerFunc) *Task {
	f := func(pc *types.PipelineContext) error {
		el, err := pc.Query().Element(selector)
		if err != nil {
			return err
		}

		if err = el.Hover(); err != nil {
			return err
		}

		if handler != nil {
			timeout := pc.Page().Timeout(time.Minute)
			wait, handle := timeout.HandleDialog()
			go handler(wait, handle)
		}

		return el.Tap()
	}
	task := &Task{op: f, desc: fmt.Sprintf("tap %q", selector)}
	return task
}

func RemoveClass(selector string, class string) *Task {
	f := func(pc *types.PipelineContext) error {
		el, err := pc.Query().Element(selector)
		if err != nil {
			return err
		}

		_, errEval := el.Eval(fmt.Sprintf(`() => this.classList.remove('%s')`, class))
		if errEval != nil {
			return errEval
		}

		return nil
	}
	task := &Task{op: f, desc: fmt.Sprintf("remove class %s from %s", class, selector)}
	return task
}

func AddClass(selector string, class string) *Task {
	f := func(pc *types.PipelineContext) error {
		el, err := pc.Query().Element(selector)
		if err != nil {
			return err
		}

		_, errEval := el.Eval(fmt.Sprintf(`() => this.classList.add('%s')`, class))
		if errEval != nil {
			return errEval
		}

		return nil
	}
	task := &Task{op: f, desc: fmt.Sprintf("add class %s to %s", class, selector)}
	return task
}

func Input(selector string, str string) *Task {
	f := func(pc *types.PipelineContext) error {
		el, err := pc.Query().Element(selector)
		if err != nil {
			return err
		}

		if err = el.SelectAllText(); err != nil {
			return err
		}

		if err = el.Input(str); err != nil {
			return err
		}

		return nil
	}
	task := &Task{op: f, desc: fmt.Sprintf("input %s to %s", str, selector)}
	return task
}

func Type(keys ...input.Key) *Task {
	f := func(pc *types.PipelineContext) error {
		if err := pc.Page().Keyboard.Type(keys...); err != nil {
			return err
		}

		return nil
	}
	task := &Task{op: f, desc: fmt.Sprintf("type keys %v", keys)}
	return task

}

func Reload() *Task {
	f := func(pc *types.PipelineContext) error {
		return pc.Page().Reload()
	}
	task := &Task{op: f, desc: "reload"}
	return task
}

func Sleep(dur time.Duration) *Task {
	f := func(pc *types.PipelineContext) error {
		time.Sleep(dur)
		return nil
	}
	task := &Task{op: f, desc: fmt.Sprintf("sleep %s", dur.String())}
	return task
}

func Stop(message string) *Task {
	f := func(pc *types.PipelineContext) error {
		return errors.New(message)
	}
	task := &Task{op: f, desc: fmt.Sprintf("stop with message %s", message)}
	return task
}

func Custom(c func(pc *types.PipelineContext) error) *Task {
	f := func(pc *types.PipelineContext) error {
		return c(pc)
	}
	task := &Task{op: f, desc: fmt.Sprintf("execute %s", util.FunctionName(c))}
	return task
}

func ForEach(selector string, ef types.EachElementFunc) *Task {
	f := func(pc *types.PipelineContext) error {
		els, err := pc.Query().Elements(selector)
		if err != nil {
			return err
		}
		for _, el := range els {
			log.Println(">>", util.FunctionName(ef))
			stop, errEl := ef(pc, el)
			if errEl != nil {
				return errEl
			} else if stop {
				break
			}
		}
		return nil
	}
	return &Task{op: f, desc: fmt.Sprintf("foreach %q", selector)}
}

func Has(selector string) types.ConditionalTask {
	f := func(pc *types.PipelineContext) (bool, error) {
		has, _, err := pc.Query().Has(selector)
		if err != nil {
			return false, err
		}

		return has, nil
	}
	return types.NewConditionalTask(f, fmt.Sprintf("has %q", selector))
}

func IsTrue(b bool) types.ConditionalTask {
	f := func(pc *types.PipelineContext) (bool, error) {
		return b, nil
	}

	return types.NewConditionalTask(f, "isTrue")
}

func ContainsText(selector, text string) types.ConditionalTask {
	f := func(pc *types.PipelineContext) (bool, error) {
		has, el, err := pc.Query().Has(selector)
		if err != nil {
			return false, err
		}
		if false == has {
			return has, nil
		}

		s, errText := el.Text()
		if errText != nil {
			return false, errText
		}

		return strings.Contains(strings.TrimSpace(s), text), nil
	}

	return types.NewConditionalTask(f, fmt.Sprintf("%q contains text %q", selector, text))
}

func Visible(selector string) types.ConditionalTask {
	f := func(pc *types.PipelineContext) (bool, error) {
		has, el, err := pc.Query().Has(selector)
		if err != nil {
			return false, err
		}
		if !has {
			return false, nil
		}

		visible, errVisible := el.Visible()
		return visible, errVisible
	}

	return types.NewConditionalTask(f, fmt.Sprintf("%q is visible", selector))
}
