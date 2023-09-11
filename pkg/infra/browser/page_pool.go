package browser

import (
	"github.com/go-rod/rod"
)

type PagePool chan *rod.Page

func NewPagePool(limit int) PagePool {
	pp := make(chan *rod.Page, limit)
	for i := 0; i < limit; i++ {
		pp <- nil
	}

	return pp
}

func (pp PagePool) Get(browser *rod.Browser, url ...string) (*rod.Page, error) {
	var err error

	p := <-pp
	if p == nil {
		p, err = newPage(browser, url...)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (pp PagePool) Put(p *rod.Page) {
	pp <- p
}

func (pp PagePool) Cleanup() error {
	for i := 0; i < cap(pp); i++ {
		p := <-pp
		if p != nil {
			if err := p.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}