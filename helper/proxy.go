package helper

import (
	"fmt"
)

type Proxy struct {
	Url string
}

// Addr -
func (p *Proxy) Addr() string {
	return fmt.Sprintf("proxy from %s", p.Url)
}

// Type -
func (p *Proxy) Type() string {
	return "proxy"
}

// String -
func (p *Proxy) String() string {
	return p.Url
}

func ParseProxy(url string) (*Proxy, error) {
	return &Proxy{
		Url: url,
	}, nil
}
