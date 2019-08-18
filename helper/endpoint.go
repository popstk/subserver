package helper

import (
	"errors"
	"strings"
)

type Endpoint interface {
	String() string
	Addr() string
	Type() string
}

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

