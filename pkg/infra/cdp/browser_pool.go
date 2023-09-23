package cdp

import (
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
)

type BrowserPoolManager interface {
	Get(settings *config.ReporterAppConfig) (*Browser, error)
	Put(b BrowserManager)
	Cleanup() error
}

type BrowserPool struct {
	Pool chan BrowserManager
}

func NewBrowserPool(setting *config.ReporterAppConfig) *BrowserPool {
	bp := make(chan BrowserManager, setting.WorkersCount)
	for i := 0; i < setting.WorkersCount; i++ {
		bp <- nil
	}

	return &BrowserPool{Pool: bp}
}

func (p *BrowserPool) Get(settings *config.ReporterAppConfig) (*Browser, error) {
	var err error

	b := <-p.Pool
	if b == nil {
		b, err = NewBrowser(settings)
		if err != nil {
			return nil, err
		}

		if err = b.Prepare(); err != nil {
			return nil, err
		}
	}

	return b.(*Browser), nil
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
