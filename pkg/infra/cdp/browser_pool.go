package cdp

type BrowserPoolManager interface {
	Get(url string) (*Browser, error)
	Put(b BrowserManager)
	Cleanup() error
}

type BrowserPool struct {
	Pool chan BrowserManager
}

func NewBrowserPool(workers int) *BrowserPool {
	bp := make(chan BrowserManager, workers)
	for i := 0; i < workers; i++ {
		bp <- nil
	}

	return &BrowserPool{Pool: bp}
}

func (p *BrowserPool) Get(url string) (*Browser, error) {
	var err error

	b := <-p.Pool
	if b == nil {
		b, err = NewBrowser(url)
		if err != nil {
			return nil, err
		}

		if err = b.Prepare(); err != nil {
			return nil, err
		}
	}

	browser, ok := b.(*Browser)
	if !ok {
		return nil, err
	}

	return browser, nil
}

func (p *BrowserPool) Put(b BrowserManager) {
	p.Pool <- b
}

func (p *BrowserPool) Cleanup() error {
	for i := 0; i < cap(p.Pool); i++ {
		b := <-p.Pool
		if b != nil {
			if err := b.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}
