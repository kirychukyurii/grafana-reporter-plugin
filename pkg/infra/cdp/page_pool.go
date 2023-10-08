package cdp

type PagePoolManager interface {
	Get(browser BrowserManager) (*Page, error)
	Put(p PageManager)
	Cleanup() error
}

type PagePool struct {
	Pool chan PageManager
}

func NewPagePool(limit int) *PagePool {
	pp := make(chan PageManager, limit)
	for i := 0; i < limit; i++ {
		pp <- nil
	}

	return &PagePool{Pool: pp}
}

func (pp PagePool) Get(browser BrowserManager) (*Page, error) {
	var err error

	p := <-pp.Pool
	if p == nil {
		p, err = NewPage(browser)
		if err != nil {
			return nil, err
		}
	}

	return p.(*Page), nil
}

func (pp PagePool) Put(p PageManager) {
	pp.Pool <- p
}

func (pp PagePool) Cleanup() error {
	for i := 0; i < cap(pp.Pool); i++ {
		p := <-pp.Pool
		if p != nil {
			if err := p.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}
