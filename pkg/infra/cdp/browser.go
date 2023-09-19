package cdp

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/dto"
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

type PageTarget struct {
	URL string
}

func NewBrowser(settings dto.ReporterAppSetting) (*Browser, error) {
	var launch string

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
