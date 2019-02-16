package main

import (
	"errors"
	"github.com/popstk/subserver/helper"
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
		lines, err := helper.Export(s.File, s.Host)
		if err != nil {
			return nil, err
		}
		return lines, nil
	}

	return nil, errors.New("Unknow type")
}
