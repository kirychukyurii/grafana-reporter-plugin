package cdp

import (
	"fmt"
	"github.com/go-rod/rod/lib/proto"
	"net"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/model"
)

type BrowserManager interface {
	Prepare() error
	Close() error

	NewPage() (*Page, error)
}

type Browser struct {
	Browser *rod.Browser
}

type PageTarget struct {
	URL string
}

func NewBrowser(settings model.ReporterAppSetting) (*Browser, error) {
	var launch string

	if settings.Browser.Url != "" {
		ips, err := net.LookupIP(settings.Browser.Url)
		if err != nil {
			return nil, fmt.Errorf("net.LookupIP: %v", err)
		}

		launch, err = launcher.ResolveURL(fmt.Sprintf("%s:9222", ips[0]))
		if err != nil {
			return nil, fmt.Errorf("launcher.ResolveURL: %v", err)
		}
	}

	return &Browser{
		Browser: rod.New().ControlURL(launch),
	}, nil
}

func (b *Browser) Prepare() error {
	if err := b.Browser.Connect(); err != nil {
		return fmt.Errorf("browser.Connect: %v", err)
	}

	return nil
}

func (b *Browser) Close() error {
	if err := b.Browser.Close(); err != nil {
		return fmt.Errorf("close browser: %v", err)
	}

	return nil
}

func (b *Browser) NewPage() (*Page, error) {
	p, err := b.Browser.Page(proto.TargetCreateTarget{URL: ""})
	if err != nil {
		return nil, fmt.Errorf("create target: %v", err)
	}

	return &Page{Page: p}, nil
}
