package task

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"

	"github.com/darimuri/rod-remote/rod_pipeline/types"
)

func Open(url string) *Task {
	f := func(pc *types.PipelineContext) error {
		return pc.Page().Navigate(url)
	}
	task := &Task{op: f}
	return task
}

func WaitLoad() *Task {
	f := func(pc *types.PipelineContext) error {
		return pc.Page().WaitLoad()
	}
	task := &Task{op: f}
	return task
}

func WaitIdle(dur time.Duration) *Task {
	f := func(pc *types.PipelineContext) error {
		return pc.Page().WaitIdle(dur)
	}
	task := &Task{op: f}
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
	task := &Task{op: f}
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
	task := &Task{op: f}
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
	task := &Task{op: f}
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
	task := &Task{op: f}
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
	task := &Task{op: f}
	return task
}

func Type(keys ...input.Key) *Task {
	f := func(pc *types.PipelineContext) error {
		if err := pc.Page().Keyboard.Type(keys...); err != nil {
			return err
		}

		return nil
	}
	task := &Task{op: f}
	return task

}

func Reload() *Task {
	f := func(pc *types.PipelineContext) error {
		return pc.Page().Reload()
	}
	task := &Task{op: f}
	return task
}

func Sleep(dur time.Duration) *Task {
	f := func(pc *types.PipelineContext) error {
		time.Sleep(dur)
		return nil
	}
	task := &Task{op: f}
	return task
}

func Stop(message string) *Task {
	f := func(pc *types.PipelineContext) error {
		return errors.New(message)
	}
	task := &Task{op: f}
	return task
}

func Custom(c func(p *rod.Page) error) *Task {
	f := func(pc *types.PipelineContext) error {
		return c(pc.Page())
	}
	task := &Task{op: f}
	return task
}

func ForEach(selector string, ef types.EachElementFunc) *Task {
	f := func(pc *types.PipelineContext) error {
		els, err := pc.Query().Elements(selector)
		if err != nil {
			return err
		}
		for _, el := range els {
			if stop, errEl := ef(el); errEl != nil {
				return errEl
			} else if stop {
				break
			}
		}
		return nil
	}
	return &Task{op: f}
}
