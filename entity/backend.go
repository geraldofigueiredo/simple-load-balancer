package entity

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	Name         string
	URL          *url.URL
	Alive        bool
	mux          sync.Mutex
	ReverseProxy *httputil.ReverseProxy
}

func NewBackend(name, urlStr string) *Backend {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return &Backend{
		Name:         name,
		URL:          u,
		Alive:        true,
		ReverseProxy: httputil.NewSingleHostReverseProxy(u),
	}
}

// SetAlive for this backend
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

// IsAlive returns true when backend is alive
func (b *Backend) IsAlive() (alive bool) {
	b.mux.Lock()
	alive = b.Alive
	b.mux.Unlock()
	return
}
