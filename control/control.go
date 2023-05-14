package control

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
)

func NewControl(b *rod.Browser) Control {
	return Control{b: b}
}

type Control struct {
	b *rod.Browser
}

func (c Control) OpenPage(url string, reuse bool) (*PageControl, error) {
	var page *rod.Page

	if reuse {
		pages, err := c.b.Pages()
		if err != nil {
			return nil, err
		}

		jsRegex := regexp.QuoteMeta(url)

		var pageNotFound *rod.ErrPageNotFound

		page, err = pages.FindByURL(jsRegex)
		if err != nil && !errors.As(err, &pageNotFound) {
			return nil, err
		}
	}

	if page == nil {
		err := rod.Try(func() {
			page = c.b.MustPage()
		})
		if err != nil {
			return nil, err
		}
	}

	if err := page.Navigate(url); err != nil {
		return nil, err
	}

	return &PageControl{Page: page}, nil
}

func (c Control) DefaultDevice(device devices.Device) {
	c.b.DefaultDevice(device)
}

type PageControl struct {
	*rod.Page
}

func (c PageControl) GetAttributesFrom(selector string, attribute string) ([]string, error) {
	var attributes []string
	err := rod.Try(func() {
		els := c.MustElements(selector)
		attributes = make([]string, 0)

		for _, el := range els {
			if attr, err := el.Attribute(attribute); err != nil {
				fmt.Println(">>>", err)
			} else if attr != nil {
				attributes = append(attributes, *attr)
			}
		}
	})

	if err != nil {
		return nil, err
	}

	return attributes, nil
}
