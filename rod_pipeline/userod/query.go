package userod

import "github.com/go-rod/rod"

func NewQuery(p *rod.Page, e *rod.Element) Query {
	return Query{p: p, e: e}
}

type Query struct {
	p *rod.Page
	e *rod.Element
}

func (q Query) Has(selector string) (bool, *rod.Element, error) {
	if q.e != nil {
		return q.e.Has(selector)
	}
	return q.p.Has(selector)
}

func (q Query) Element(selector string) (*rod.Element, error) {
	if q.e != nil {
		return q.e.Element(selector)
	}
	return q.p.Element(selector)

}

func (q Query) Elements(selector string) (rod.Elements, error) {
	if q.e != nil {
		return q.e.Elements(selector)
	}
	return q.p.Elements(selector)
}
