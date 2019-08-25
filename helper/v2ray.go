package helper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"strings"
)

// VmessURL implement interface Endpoint
// from https://github.com/2dust/v2rayN/wiki/分享链接格式说明(ver-2)
type VmessURL struct {
	Version  Number `json:"v"`
	AlterID  Number `json:"aid"`
	Ps       string `json:"ps"`
	Port     string `json:"port"`
	ID       string `json:"id"`
	Network  string `json:"net"`
	Security string `json:"tls"`
	Add      string `json:"add"`
	FakePath string `json:"path"`
	FakeType string `json:"type"`
	FakeHost string `json:"host"`
}

// NewVmessURL -
func NewVmessURL(u string) (*VmessURL, error) {
	const prefix = "vmess://"

	if !strings.HasPrefix(u, prefix) {
		return nil, errors.New("helper: invalid vmess url")
	}

	var vu VmessURL
	data, err := Base64Decode(u[len(prefix):])
	if err != nil {
		log.Println(u)
		return nil, errors.Wrap(err, "helper: can not decode data")
	}

	if err = json.Unmarshal([]byte(data), &vu); err != nil {
		return nil, errors.Wrap(err, "helper: can not unmarshal data")
	}

	return &vu, nil
}

// Addr -
func (v *VmessURL) Addr() string {
	return fmt.Sprintf("%s:%s", v.Add, v.Port)
}

// Type -
func (v *VmessURL) Type() string {
	return "vmess"
}

// String -
func (v *VmessURL) String() string {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return ""
	}

	return "vmess://" + base64.StdEncoding.EncodeToString(data)
}

// VmessParser -
var VmessParser = JSONParser{
	Filed: map[string]FieldParser{
		"protocol":      JSONPathHandler("protocol"),
		"port":          JSONPathHandler("port"),
		"id":            JSONPathHandler("settings.clients.0.id"),
		"alterId":       JSONPathHandler("settings.clients.0.alterId"),
		"network":       JSONPathHandler("streamSettings.network"),
		"security":      JSONPathHandler("streamSettings.security"),
		"http.host":     JSONPathHandler("streamSettings.httpSettings.host.0"),
		"http.path":     JSONPathHandler("streamSettings.httpSettings.path"),
		"ws.host":       JSONPathHandler("streamSettings.wsSettings.headers.Host"),
		"ws.path":       JSONPathHandler("streamSettings.wsSettings.path"),
		"kcp.type":      JSONPathHandler("streamSettings.kcpSettings.header.type"),
		"quic.type":     JSONPathHandler("streamSettings.quicSettings.header.type"),
		"quic.security": JSONPathHandler("streamSettings.quicSettings.security"),
		"quic.key":      JSONPathHandler("streamSettings.quicSettings.key"),
		"servername":    JSONPathHandler("tlsSettings.serverName"),
	},
	DefaultField: map[string]string{
		"add": "",
	},
	PostHandler: func(m map[string]string, tag string) (Endpoint, error) {
		alterID, err := strconv.Atoi(m["alterId"])
		if err != nil {
			return nil, errors.Wrap(err, "VmessParser")
		}

		vu := VmessURL{
			Ps:       tag,
			Port:     m["port"],
			ID:       m["id"],
			AlterID:  Number(alterID),
			Network:  m["network"],
			Security: m["security"],
			Add:      m["add"],
			Version:  2,
		}

		if m["network"] == "http" {
			vu.FakeHost = m["http.host"]
			vu.FakePath = m["http.path"]
			vu.Network = "h2"
		} else if m["network"] == "ws" {
			vu.FakeHost = m["ws.host"]
			vu.FakePath = m["ws.path"]
		} else if m["network"] == "kcp" {
			if m["kcp.type"] == "" {
				m["kcp.type"] = "none"
			}
			vu.FakeType = m["kcp.type"]
		} else if m["network"] == "quic" {
			if m["quic.type"] == "" {
				m["quic.type"] = "none"
			}
			vu.FakeType = m["quic.type"]
			vu.FakeHost = m["quic.security"]
			vu.FakePath = m["quic.key"]
		}

		if vu.Port == "443" {
			vu.Security = "tls"
		}

		return &vu, nil
	},
}
