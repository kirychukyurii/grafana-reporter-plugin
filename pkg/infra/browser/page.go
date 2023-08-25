package browser

import (
	"fmt"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
)

func (b *browser) Page(url ...string) (*rod.Page, error) {
	var err error

	b.page, err = b.browser.Page(proto.TargetCreateTarget{
		URL: strings.Join(url, "/"),
	})

	if err != nil {
		return nil, fmt.Errorf("open page: %v", err)
	}

	// Start to analyze request events
	wait := b.page.MustWaitRequestIdle()

	// Wait until there's no active requests
	wait()

	return b.page, nil
}

func (b *browser) ExtraHeaders(headers []string) (func(), error) {
	f, err := b.page.SetExtraHeaders(headers)
	if err != nil {
		return nil, fmt.Errorf("set extra headers: %v", err)
	}

	return f, nil
}

func (b *browser) Click(selector string, jsRegex string) error {
	e, err := b.page.ElementR(selector, jsRegex)
	if err != nil {
		return fmt.Errorf("read element: %v", err)
	}

	if err = e.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("click: %v", err)
	}

	return nil
}

func (b *browser) Screenshot() {
	//TODO implement me
	panic("implement me")
}

func (b *browser) DownloadFile(path string) error {
	wait := b.browser.MustWaitDownload()
	if err := utils.OutputFile(path, wait()); err != nil {
		return fmt.Errorf("save file: %v", err)
	}

	return nil
}
