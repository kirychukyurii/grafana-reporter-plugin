package browser

import (
	"github.com/go-rod/rod"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/model"
)

type BrowserPool chan *rod.Browser

func NewBrowserPool(limit int) BrowserPool {
	bp := make(chan *rod.Browser, limit)
	for i := 0; i < limit; i++ {
		bp <- nil
	}

	return bp
}

func (p BrowserPool) Get(settings model.ReporterAppSetting) (*rod.Browser, error) {
	var err error

	b := <-p
	if b == nil {
		b, err = newBrowser(settings)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

func (p BrowserPool) Put(b *rod.Browser) {
	p <- b
}

func (p BrowserPool) Cleanup() error {
	for i := 0; i < cap(p); i++ {
		b := <-p
		if b != nil {
			if err := b.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}
