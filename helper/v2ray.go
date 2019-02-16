package helper

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type v2ray struct {
	Inbounds []inboundObject `json:"inbounds"`
}

type httpSettings struct {
	Host []string `json:"host"`
	Path string   `json:"path"`
}

type wsSettings struct {
	Headers map[string]string `json:"headers"`
	Path    string            `json:"path"`
}

type quicHeader struct {
	Type string `json:"type"`
}

type quicSettings struct {
	Security string     `json:"security"`
	Key      string     `json:"key"`
	Header   quicHeader `json:"header"`
}

type inboundStream struct {
	Network  string       `json:"network"`
	Security string       `json:"security"`
	Http     httpSettings `json:"httpSettings"`
	Ws       wsSettings   `json:"wsSettings"`
	Quic     quicSettings `json:"quicSettings"`
}

type inboundObject struct {
	Port           int             `json:"port"`
	Listen         string          `json:"listen"`
	Protocol       string          `json:"protocol"`
	Settings       json.RawMessage `json:"settings"`
	StreamSettings inboundStream   `json:"streamSettings"`
}

type inboundVmessClient struct {
	ID      string `json:"id"`
	AlterID int    `json:"alterId"`
	Email   string `json:"email"`
}

type inboundVmessDefault struct {
	AlterID string `json:"alterId"`
}

type inboundVmess struct {
	Clients []inboundVmessClient `json:"clients"`
	Default inboundVmessDefault  `json:"default"`
}

type inboundSS struct {
	Email    string `json:"email"`
	Method   string `json:"method"`
	Password string `json:"password"`
	Ota      string `json:"ota"`
	Network  string `json:"network"`
}

// generateVmess -
// from https://github.com/2dust/v2rayN/wiki/分享链接格式说明(ver-2)
func generateVmess(i inboundObject) ([]string, error) {
	account := make([]string, 0)

	var vmess inboundVmess
	if err := json.Unmarshal([]byte(i.Settings), &vmess); err != nil {
		return nil, err
	}

	for _, client := range vmess.Clients {
		conf := map[string]string{
			"v":    "2",
			"add":  i.Listen,
			"port": strconv.Itoa(i.Port),
			"id":   client.ID,
			"aid":  strconv.Itoa(client.AlterID),
		}

		if i.StreamSettings.Network != "" {
			conf["net"] = i.StreamSettings.Network
		}

		if i.StreamSettings.Security != "" {
			conf["tls"] = i.StreamSettings.Security
		}

		data, err := json.Marshal(conf)
		if err != nil {
			return nil, err
		}

		account = append(account, "vmess://"+base64.StdEncoding.EncodeToString(data))
	}

	return account, nil
}

func generateSS(i inboundObject) ([]string, error) {
	account := make([]string, 0)
	var ss inboundSS
	if err := json.Unmarshal([]byte(i.Settings), &ss); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("%s:%s@%s:%d", ss.Method, ss.Password, i.Listen, i.Port)
	account = append(account, "ss://"+base64.StdEncoding.EncodeToString([]byte(u)))
	return account, nil
}

// Export -
func Export(filepath, host string) ([]string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	var conf v2ray
	reader := bufio.NewReader(f)
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&conf)
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0)
	for _, inbound := range conf.Inbounds {
		inbound.Listen = host // replace host

		if inbound.Protocol == "vmess" {
			link, err := generateVmess(inbound)
			if err != nil {
				return nil, err
			}
			ret = append(ret, link...)
		} else if inbound.Protocol == "shadowsocks" {
			link, err := generateSS(inbound)
			if err != nil {
				return nil, err
			}
			ret = append(ret, link...)
		}
	}

	return ret, nil
}
