package browser

import (
	"fmt"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type Pager interface {
	Get(browser *rod.Browser) (*rod.Page, error)
	Put(p *rod.Page)
	Cleanup() error
}

func newPage(browser *rod.Browser, url ...string) (*rod.Page, error) {
	p, err := browser.Page(proto.TargetCreateTarget{
		URL: strings.Join(url, "/"),
	})

	if err != nil {
		return nil, err
	}

	return p, nil
}

type page struct {
	page *rod.Page
}

func (p *page) New(browser *rod.Browser, url ...string) (*rod.Page, error) {
	var err error

	p.page, err = browser.Page(proto.TargetCreateTarget{
		URL: strings.Join(url, "/"),
	})

	if err != nil {
		return nil, fmt.Errorf("open page: %v", err)
	}

	// Start to analyze request events
	wait := p.page.MustWaitRequestIdle()

	// Wait until there's no active requests
	wait()

	return p.page, nil
}

func (p *page) ExtraHeaders(headers []string) (func(), error) {
	f, err := p.page.SetExtraHeaders(headers)
	if err != nil {
		return nil, fmt.Errorf("set extra headers: %v", err)
	}

	return f, nil
}

func (p *page) Click(selector string, jsRegex string) error {
	e, err := p.page.ElementR(selector, jsRegex)
	if err != nil {
		return fmt.Errorf("read element: %v", err)
	}

	if err = e.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("click: %v", err)
	}

	return nil
}

func (p *page) Screenshot() {
	//TODO implement me
	panic("implement me")
}

/*
func (p *page) DownloadFile(path string) error {
	wait := p.browser.MustWaitDownload()
	if err := utils.OutputFile(path, wait()); err != nil {
		return fmt.Errorf("save file: %v", err)
	}

	return nil
}

*/
