package cdp

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type PageManager interface {
	Prepare(url string, headers []string, viewport *PageViewportOpts) (*Page, error)
	Close() error
}

type Page struct {
	Page *rod.Page
}

type PageViewportOpts struct {
	Width  int
	Height int
}

func NewPage(b BrowserManager) (*Page, error) {
	p, err := b.NewPage()

	if err != nil {
		return nil, fmt.Errorf("new page: %v", err)
	}

	return p, nil
}

func (p *Page) Prepare(url string, headers []string, viewport *PageViewportOpts) (*Page, error) {
	if viewport != nil {
		if err := p.Page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{Width: viewport.Width, Height: viewport.Height}); err != nil {
			return nil, err
		}
	}

	if len(headers) > 0 {
		_, err := p.Page.SetExtraHeaders(headers)
		if err != nil {
			return nil, err
		}
	}

	if err := p.Page.Navigate(url); err != nil {
		return nil, fmt.Errorf("navigate to url: %v", err)
	}

	if err := p.Page.WaitLoad(); err != nil {
		return nil, fmt.Errorf("wait load: %v", err)
	}

	w := p.Page.MustWaitRequestIdle()
	w()

	return p, nil
}

func (p *Page) Close() error {
	if err := p.Page.Close(); err != nil {
		return fmt.Errorf("close page: %v", err)
	}

	return nil
}
