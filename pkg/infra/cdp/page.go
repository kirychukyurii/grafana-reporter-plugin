package cdp

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/util"
)

type PageManager interface {
	Prepare(url string, headers []string, viewport *PageViewportOpts) error
	Close() error

	Eval(s string, args ...interface{}) (*proto.RuntimeRemoteObject, error)
	Scroll(offsetX, offsetY float64, steps int) error
	ScrollHeight() (float64, float64, error)
	ScrollDown(sleep time.Duration) error
	Elements(selector string) (rod.Elements, error)
	Screenshot(file string, full bool) error
	ScreenshotFullPage(file string) error
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

func (p *Page) Prepare(url string, headers []string, viewport *PageViewportOpts) error {
	p.Page.StopLoading()
	if viewport != nil {
		if err := p.Page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{Width: viewport.Width, Height: viewport.Height}); err != nil {
			return fmt.Errorf("set viewport: %v", err)
		}
	}

	if len(headers) > 0 {
		_, err := p.Page.SetExtraHeaders(headers)
		if err != nil {
			return err
		}
	}

	if err := p.Page.Navigate(url); err != nil {
		return fmt.Errorf("navigate to url: %v", err)
	}

	if err := p.Page.WaitLoad(); err != nil {
		return fmt.Errorf("wait load: %v", err)
	}

	w := p.Page.WaitRequestIdle(200*time.Millisecond, nil, []string{}, nil)
	w()

	return nil
}

func (p *Page) Close() error {
	if err := p.Page.Close(); err != nil {
		return fmt.Errorf("close page: %v", err)
	}

	return nil
}

func (p *Page) Eval(s string, args ...interface{}) (*proto.RuntimeRemoteObject, error) {
	obj, err := p.Page.Eval(s, args)
	if err != nil {
		return nil, fmt.Errorf("evaluate js: %v", err)
	}

	return obj, nil
}

func (p *Page) Scroll(offsetX, offsetY float64, steps int) error {
	if err := p.Page.Mouse.Scroll(offsetX, offsetY, steps); err != nil {
		return fmt.Errorf("scroll: %v", err)
	}

	return nil
}

func (p *Page) ScrollHeight() (float64, float64, error) {
	scrollHeightObj, err := p.Eval("() => document.documentElement.scrollHeight")
	if err != nil {
		return 0, 0, fmt.Errorf("scrollHeightObj: %v", err)
	}

	clientHeightObj, err := p.Eval("() => document.documentElement.clientHeight")
	if err != nil {
		return 0, 0, fmt.Errorf("clientHeightObj: %v", err)
	}

	return scrollHeightObj.Value.Num(), clientHeightObj.Value.Num(), nil
}

func (p *Page) ScrollDown(sleep time.Duration) error {
	scrollHeight, clientHeight, err := p.ScrollHeight()
	if err != nil {
		return err
	}

	if scrollHeight < clientHeight {
	}

	scrolls := int(scrollHeight / clientHeight)
	for i := 1; i < scrolls; i++ {
		if err = p.Scroll(0, clientHeight, 0); err != nil {
			return fmt.Errorf("scroll: %v", err)
		}

		time.Sleep(sleep * time.Millisecond)
	}

	if err = p.Scroll(0, 0, 0); err != nil {
		return fmt.Errorf("scroll to 0,0: %v", err)
	}

	return nil
}

func (p *Page) Elements(selector string) (rod.Elements, error) {
	e, err := p.Page.Elements(selector)
	if err != nil {
		return nil, fmt.Errorf("elements: %v", err)
	}

	return e, err
}

func (p *Page) Screenshot(file string, full bool) error {
	since := time.Now()
	defer func() { backend.Logger.Debug(util.TimeTrack(since)) }()

	bin, err := p.Page.Screenshot(full, nil)
	if err != nil {
		return fmt.Errorf("screenshot: %v", err)
	}

	if err = OutputFile(file, bin); err != nil {
		return err
	}

	return nil
}

func (p *Page) ScreenshotFullPage(file string) error {
	if err := p.ScrollDown(500); err != nil {
		return err
	}

	if err := p.Screenshot(file, true); err != nil {
		return err
	}

	return nil
}

func (p *Page) Element(selector, jsRegex string) (*rod.Element, error) {
	var (
		e   *rod.Element
		err error
	)

	if jsRegex != "" {
		e, err = p.Page.ElementR(selector, jsRegex)
		if err != nil {
			return nil, err
		}
	} else {
		e, err = p.Page.Element(selector)
		if err != nil {
			return nil, err
		}
	}

	return e, nil
}
