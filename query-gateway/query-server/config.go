package query_server

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/popstk/subserver/helper"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// Config -
type Config struct {
	Address string              `json:"addr"`
	Valid   map[string][]Source `json:"valid"`
}

// Fmt -
type Fmt struct {
	Vmess string `json:"vmess"`
	SS string  `json:"ss"`
}

// Source -
type Source struct {
	Type string `json:"type"`
	File string `json:"file"`
	Host string `json:"host"`
	Addr string `json:"addr"`
	Fmt Fmt `json:"fmt"`
	Rewrite map[string]string `json:"rewrite"`
}

const (
	defaultVmessName = "{protocol}-{network}"
	defaultSSName    = "{protocol}"
)

// Parse -
func (s *Source) Parse() ([]string, error) {
	if s.Type == "raw" {
		lines, err := helper.ReadLines(s.File)
		if err != nil {
			return nil, err
		}
		return lines, nil
	}

	if s.Type == "v2ray" {
		lines, err := Export(s)
		if err != nil {
			return nil, err
		}
		return lines, nil
	}

	if s.Type == "url" {
		rsp, err := http.Get(s.Addr)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()

		data, err :=ioutil.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		data, err  = base64.StdEncoding.DecodeString(string(data))
		if err != nil {
			return nil, err
		}

		return strings.Split(string(data), "\n"), nil

	}

	return nil, errors.New("Unknow type: "+s.Type)
}

// Export -
func Export(s *Source) ([]string, error) {
	data, err := ioutil.ReadFile(s.File)
	if err !=nil {
		return nil, err
	}

	if s.Fmt.Vmess == "" {
		s.Fmt.Vmess = defaultVmessName
	}

	if s.Fmt.SS == "" {
		s.Fmt.SS = defaultSSName
	}

	value := gjson.Get(string(data), "inbounds")
	if !value.IsArray() {
		return nil, errors.New("unknown format")
	}

	ret := make([]string, 0)

	for _, inbound := range value.Array() {
		host := s.Host
		port := inbound.Get("port").String()
		listen := inbound.Get("listen").String()

		add, exist := s.Rewrite[port]
		if exist {
			host, port, err = net.SplitHostPort(add)
			if err != nil {
				return nil, err
			}
		} else if listen == "127.0.0.1" || listen == "localhost" {
			return nil, errors.New(fmt.Sprint("loopback inbound without rewrite: ", port))
		}

		protocol := inbound.Get("protocol").String()
		if protocol == "vmess" {
			helper.VmessParser.TagFmt = s.Fmt.Vmess
			helper.VmessParser.DefaultField["add"] = host
			helper.VmessParser.DefaultField["port"] = port
			u, err := helper.VmessParser.Parse(inbound)
			if err != nil {
				return nil, err
			}
			ret = append(ret, u)
		} else if protocol == "shadowsocks" {
			helper.SSParser.TagFmt = s.Fmt.SS
			helper.SSParser.DefaultField["host"] = host
			helper.SSParser.DefaultField["port"] = port
			u, err := helper.SSParser.Parse(inbound)
			if err != nil {
				return nil, err
			}
			ret = append(ret, u)
		}
	}


	return ret, nil
}