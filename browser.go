package main

import (
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
	"io/ioutil"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/ysmood/gson"
)

type Headless struct {
	launcher *launcher.Launcher
	browser  *rod.Browser
	page     *rod.Page
}

func NewHeadlessBrowser() Headless {
	l := launcher.New().
		Headless(false).
		Devtools(true)

	url := l.MustLaunch()
	browser := rod.New().ControlURL(url).Trace(true).
		SlowMotion(2 * time.Second).MustConnect()

	//launcher.Open(browser.ServeMonitor(""))

	return Headless{
		launcher: l,
		browser:  browser,
	}
}

func (h *Headless) Page(url string) {
	h.page = h.browser.MustPage(url)
	wait := h.page.MustWaitRequestIdle()
	wait()

}

func (h *Headless) SetHeaders(headers []string) error {
	_, err := h.page.SetExtraHeaders(headers)
	if err != nil {
		return err
	}

	return nil
}

func (h *Headless) Screenshot(path string) error {
	img, _ := h.page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format:  proto.PageCaptureScreenshotFormatJpeg,
		Quality: gson.Int(90),
		Clip: &proto.PageViewport{
			X:      0,
			Y:      0,
			Width:  300,
			Height: 200,
			Scale:  1,
		},
		FromSurface: true,
	})

	if err := utils.OutputFile(path, img); err != nil {
		return err
	}

	return nil
}

func main() {
	var h Headless

	//	page := browser.MustPage("https://webitel.cashx.lk/grafana/d/RVJgUcnnk/agents?orgId=1&inspect=4&inspectTab=data").MustWaitLoad()
	_ = h.browser.MustPage("https://webhook.site/9e65782a-732f-4809-aa23-417eb8e830a1").MustWaitLoad()

	// _, err := page.SetExtraHeaders([]string{"Authorization", "Bearer glsa_1gIxrVgQBAWbrJirrpewkfa0X1dt6d5X_104c9da3"})
	_, err := h.page.SetExtraHeaders([]string{"Authorization", "Basic YWRtaW46ME1QMVFjSm5rQTlC"})
	if err != nil {
		return
	}

	// Start to analyze request events
	wait := h.page.MustWaitRequestIdle()

	// Wait until there's no active requests
	wait()
	// simple version
	h.page.MustScreenshot("my.png")

	// customization version
	img, _ := h.page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format:  proto.PageCaptureScreenshotFormatJpeg,
		Quality: gson.Int(90),
		Clip: &proto.PageViewport{
			X:      0,
			Y:      0,
			Width:  300,
			Height: 200,
			Scale:  1,
		},
		FromSurface: true,
	})

	err = ioutil.WriteFile("fullScreenshot.png", img, 0o644)
	if err != nil {
		panic(err)
	}

	dwait := h.browser.MustWaitDownload()
	h.page.MustElementR("span", "Download CSV").MustClick()

	_ = utils.OutputFile("t.csv", dwait())
}
