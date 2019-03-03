package main

import (
	"encoding/base64"
	"errors"
	"github.com/popstk/subserver/helper"
	"io/ioutil"
	"net/http"
	"strings"
)

// Config -
type Config struct {
	Address string              `json:"addr"`
	Valid   map[string][]Source `json:"valid"`
}

// Source -
type Source struct {
	Type string `json:"type"`
	File string `json:"file"`
	Host string `json:"host"`
	Addr string `json:"addr"`
	VmessFmt string `json:"vmess-fmt"`
	SSFmt string  `json:"ss-fmt"`
}

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
		lines, err := helper.Export(s.File,s.Host, s.VmessFmt, s.SSFmt)
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
