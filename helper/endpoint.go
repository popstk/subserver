package helper

import (
	"errors"
	"strings"
)

// Endpoint endpoint
type Endpoint interface {
	String() string
	Addr() string
	Type() string
}

// ParseEndpoint parse endpoint by uri proto
func ParseEndpoint(u string) (Endpoint, error) {
	u = strings.TrimSpace(u)

	if strings.HasPrefix(u, "vmess") {
		return NewVmessURL(u)
	}

	if strings.HasPrefix(u, "ss") {
		return NewShadowsocksURL(u)
	}

	return nil, errors.New("unknown protocol")
}
