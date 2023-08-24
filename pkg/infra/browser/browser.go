package browser

import (
	"context"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/cdp"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"net"
	"time"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/models"
)

type Browser interface {
}

type browser struct {
	browser *rod.Browser
}

func New(settings models.ReporterAppSetting) Browser {
	browserLoaded := browserTimeoutDetector(10 * time.Second)
	defer browserLoaded()

	launch := launcher.New().
		Headless(true).
		Leakless(true).
		Devtools(false).
		NoSandbox(true).
		Set("disable-web-security") // TODO: ensure we have proper CORS

	if settings.Browser.Url != "" {
		host, port, _ := net.SplitHostPort(settings.Browser.Url)
		launch = launch.Set("remote-debugging-address", host).Set(flags.RemoteDebuggingPort, port)
	}

	if settings.Browser.BinPath != "" {
		launch = launch.Bin(settings.Browser.BinPath)
	}

	defer func() {
		launch.Kill()
		avoidStall(3*time.Second, launch.Cleanup)
	}()

	url, err := launch.Launch()
	client := cdp.New()

	b := rod.New().
		Timeout(time.Minute).
		Sleeper(MaxDuration(5 * time.Second)).
		Client(client).
		Context(ctx)

	defer ctx.Check(func() error {
		// browser.Close may sometimes return context.Canceled.
		return errs2.IgnoreCanceled(b.Close())
	})

	return browser{}
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
