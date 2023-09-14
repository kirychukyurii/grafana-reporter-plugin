package cdp

import (
	"github.com/go-rod/rod"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/model"
)

type BrowserPoolManager interface {
	Get(settings model.ReporterAppSetting) (*rod.Browser, error)
	Put(b BrowserManager)
	Cleanup() error
}

type BrowserPool struct {
	Pool chan BrowserManager
}

func NewBrowserPool(limit int) *BrowserPool {
	bp := make(chan BrowserManager, limit)
	for i := 0; i < limit; i++ {
		bp <- nil
	}

	return &BrowserPool{Pool: bp}
}

func (p *BrowserPool) Get(settings model.ReporterAppSetting) (*Browser, error) {
	var (
		browser *Browser
		err     error
	)

	b := <-p.Pool
	if b == nil {
		browser, err = NewBrowser(settings)
		if err != nil {
			return nil, err
		}
	}

	return browser, nil
}

func (p *BrowserPool) Put(b BrowserManager) {
	p.Pool <- b
}

func (p *BrowserPool) Cleanup() error {
	for i := 0; i < cap(p.Pool); i++ {
		b := <-p.Pool
		if b != nil {
			if err := b.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}
