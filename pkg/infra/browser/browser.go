package browser

import (
	"context"
	"fmt"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/model"
	"net"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Browser interface {
	Get(settings model.ReporterAppSetting) (*rod.Browser, error)
	Put(b *rod.Browser)
	Cleanup() error
}

func newBrowser(settings model.ReporterAppSetting) (*rod.Browser, error) {
	ips, err := net.LookupIP("chrome")
	if err != nil {
		return nil, fmt.Errorf("net.LookupIP: %v", err)
	}

	launch, err := launcher.ResolveURL(fmt.Sprintf("%s:9222", ips[0]))
	if err != nil {
		return nil, fmt.Errorf("launcher.ResolveURL: %v", err)
	}

	b := rod.New().ControlURL(launch)
	if err = b.Connect(); err != nil {
		return nil, fmt.Errorf("browser.Connect [%s]: %v", launch, err)
	}

	//browserLoaded := browserTimeoutDetector(10 * time.Second)
	//defer browserLoaded()

	/*
		launch := launcher.New().
			Headless(true).
			Leakless(true).
			Devtools(false).
			NoSandbox(true).
			Set("disable-web-security"). // TODO: ensure we have proper CORS

		if settings.Browser.Url != "" {
			host, port, _ := net.SplitHostPort(settings.Browser.Url)
			launch = launch.Set("remote-debugging-address", host).Set(flags.RemoteDebuggingPort, port)
		}

	*/

	/*
		if settings.Browser.BinPath != "" {
			launch = launch.Bin(settings.Browser.BinPath)
		}

		defer func() {
			launch.Kill()
			avoidStall(3*time.Second, launch.Cleanup)
		}()



		url, err := launch.Launch()
		if err != nil {
			return nil, fmt.Errorf("browser launcher: %v", err)
		}

		b := rod.New().Timeout(time.Minute).ControlURL(url)
	*/

	return b, nil
}

func browserTimeoutDetector(duration time.Duration) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		t := time.NewTimer(duration)
		defer t.Stop()
		select {
		case <-t.C:
			panic("timeout for starting browser exceeded")
		case <-ctx.Done():
			return
		}
	}()
	return cancel
}

func avoidStall(maxDuration time.Duration, fn func()) {
	done := make(chan struct{})
	go func() {
		fn()
		close(done)
	}()

	timeout := time.NewTicker(maxDuration)
	defer timeout.Stop()
	select {
	case <-done:
	case <-timeout.C:
		fmt.Printf("go-rod did not shutdown within %v\n", maxDuration)
	}
}
