package cdp

import (
	"fmt"
	"net"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
)

type BrowserManager interface {
	Prepare() error
	Close() error

	NewPage() (*Page, error)
	WaitDownload(dir string) func() (info *proto.PageDownloadWillBegin)
}

type Browser struct {
	Browser *rod.Browser
}

func NewBrowser(settings *config.ReporterAppConfig) (*Browser, error) {
	var launch string

	if settings.BrowserConfig.URL != "" {
		ips, err := net.LookupIP(settings.BrowserConfig.URL)
		if err != nil {
			return nil, fmt.Errorf("net.LookupIP: %v", err)
		}

		launch, err = launcher.ResolveURL(fmt.Sprintf("%s:9222", ips[0]))
		if err != nil {
			return nil, fmt.Errorf("launcher.ResolveURL: %v", err)
		}
	}

	browser := rod.New().ControlURL(launch).Logger(log.New()).Trace(true)

	return &Browser{
		Browser: browser,
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

func (b *Browser) WaitDownload(dir string) func() (info *proto.PageDownloadWillBegin) {
	f := b.Browser.WaitDownload(dir)

	return f
}
