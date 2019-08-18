package helper

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

type ShadowsocksURL struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Method   string `json:"method"`
	Password string `json:"password"`
	Protocol string `json:"protocol"`
	Tag      string `json:"tag"`
}

func NewShadowsocksURL(u string) (*ShadowsocksURL, error) {
	const prefix = "ss://"
	if strings.HasPrefix(u, prefix) {
		return nil, errors.New("helper: invalid shadowsocks url")
	}

	var su ShadowsocksURL
	p := strings.SplitN(u[len(prefix):], "#", 2)
	if len(p) == 2 {
		su.Tag = p[1]
	}

	data, err := base64.StdEncoding.DecodeString(p[0])
	if err != nil {
		return nil, errors.Wrap(err, "helper: can not decode url")
	}

	uri, err := url.Parse("ss://" + string(data))
	if err != nil {
		return nil, errors.Wrap(err, "helper: invalid url")
	}

	pwd, _ := uri.User.Password()
	su.Method = uri.User.Username()
	su.Password = pwd
	su.Host = uri.Hostname()
	su.Port = uri.Port()

	return &su, nil
}

func (s *ShadowsocksURL) Addr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

func (s *ShadowsocksURL) Type() string {
	return "shadowsocks"
}

func (s *ShadowsocksURL) String() string {
	data := fmt.Sprintf("%s:%s@%s:%s", s.Method, s.Password, s.Host, s.Port)
	return fmt.Sprintf("ss://%s#%s", base64.StdEncoding.EncodeToString([]byte(data)), s.Tag)
}

var SSParser = JSONParser{
	Filed: map[string]FieldParser{
		"protocol": JSONPathHandler("protocol"),
		"port":     JSONPathHandler("port"),
		"method":   JSONPathHandler("settings.method"),
		"password": JSONPathHandler("settings.password"),
	},
	DefaultField: map[string]string{
		"host": "",
	},
	PostHandler: func(m map[string]string, tag string) (Endpoint, error) {
		su := ShadowsocksURL{
			Host:     m["host"],
			Port:     m["port"],
			Method:   m["method"],
			Password: m["password"],
			Tag:      tag,
		}

		return &su, nil
	},
}
