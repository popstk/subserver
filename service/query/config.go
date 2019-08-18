package query

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/popstk/subserver/helper"
	"github.com/tidwall/gjson"
)

// Config -
type Config struct {
	Address string              `json:"addr"`
	Valid   map[string][]Source `json:"valid"`
}

// Fmt -
type Fmt struct {
	Vmess string `json:"vmess"`
	SS    string `json:"ss"`
}

func LoadConfig() error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	configMutex.Lock()
	defer configMutex.Unlock()

	if err = json.Unmarshal(data, &config); err != nil {
		return err
	}

	log.Println("config updated")
	return nil
}

func ListenConfig(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err = watcher.Add(configFile); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case ev := <-watcher.Events:
			if ev.Op == fsnotify.Write {
				_ = LoadConfig()
			}
		}
	}
}

type Source struct {
	Type    string            `json:"type"`
	File    string            `json:"file"`
	Host    string            `json:"host"`
	Addr    string            `json:"addr"`
	Fmt     Fmt               `json:"fmt"`
	Rewrite map[string]string `json:"rewrite"`
}

const (
	defaultVmessName = "{protocol}-{network}"
	defaultSSName    = "{protocol}"
)

func ParseV2rayConfig(s *Source) ([]helper.Endpoint, error) {
	data, err := ioutil.ReadFile(s.File)
	if err != nil {
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

	ret := make([]helper.Endpoint, 0)

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

func URLToEndpoint(uris ...string) []helper.Endpoint {
	var nodes []helper.Endpoint
	for _, u := range uris {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}

		node, err := helper.ParseEndpoint(u)
		if err != nil {
			log.Println(err, " => ", u)
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func RawNodeHandler(s *Source) ([]helper.Endpoint, error) {
	uris, err := helper.ReadLines(s.File)
	if err != nil {
		return nil, err
	}

	return URLToEndpoint(uris...), nil
}

func V2rayNodeHandler(s *Source) ([]helper.Endpoint, error) {
	return ParseV2rayConfig(s)
}

func UrlNodeHandler(s *Source) ([]helper.Endpoint, error) {
	if s.Addr == "" {
		return nil, errors.New("addr: url is empty")
	}

	rsp, err := http.Get(s.Addr)
	if err != nil {
		return nil, errors.Wrap(err, "UrlNodeHandler")
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "UrlNodeHandler")
	}
	_ = rsp.Body.Close()


	nodes, err := helper.Base64Decode(string(data))
	if err != nil {
		return nil, errors.Wrap(err, "UrlNodeHandler")
	}

	urls := strings.Split(nodes, "\n")
	return URLToEndpoint(urls...), nil
}

var handlerMap = map[string]func(s *Source) ([]helper.Endpoint, error){
	"raw":   RawNodeHandler,
	"v2ray": V2rayNodeHandler,
	"url":   UrlNodeHandler,
}

func ParseConfig(uuid string) ([]helper.Endpoint, error) {
	nodes := make([]helper.Endpoint, 0)
	m := make(map[string]bool)
	q := []string{uuid}

	configMutex.Lock()
	defer configMutex.Unlock()

	for len(q) > 0 {
		uuid := q[0]
		q = q[1:]
		m[uuid] = true

		array, exist := config.Valid[uuid]
		if !exist {
			return nil, errors.New("Invalid uuid: " + uuid)
		}

		for _, s := range array {
			if s.Type == "sub" {
				_, exist := m[s.Addr]
				if exist {
					q = append(q, s.Addr)
				}
				continue
			}

			handler, exist := handlerMap[s.Type]
			if !exist {
				log.Println("invalid type: ", s.Type)
			}

			eps, err := handler(&s)
			if err != nil {
				log.Println(s.Type, " => ", err)
				continue
			}

			nodes = append(nodes, eps...)
		}
	}

	return nodes, nil
}
