package helper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// from https://github.com/2dust/v2rayN/wiki/分享链接格式说明(ver-2)
var VmessParser = JSONParser{
	Filed: map[string]FieldParser{
		"protocol":     JSONPathHandler("protocol"),
		"port": JSONPathHandler("port"),
		"id":   JSONPathHandler("settings.clients.0.id"),
		"alterId":  JSONPathHandler("settings.clients.0.alterId"),
		"network":  JSONPathHandler("streamSettings.network"),
		"security":  JSONPathHandler("streamSettings.security"),
		"http.host": JSONPathHandler("streamSettings.httpSettings.host.0"),
		"http.path": JSONPathHandler("streamSettings.httpSettings.path"),
		"ws.host": JSONPathHandler("streamSettings.wsSettings.headers.Host"),
		"ws.path": JSONPathHandler("streamSettings.wsSettings.path"),
		"kcp.type": JSONPathHandler("streamSettings.kcpSettings.header.type"),
		"quic.type": JSONPathHandler("streamSettings.quicSettings.header.type"),
		"servername": JSONPathHandler("tlsSettings.serverName"),
	},
	DefaultField: map[string]string{
		"add": "",
	},
	PostHandler: func(m map[string]string, tag string) (string, error) {
		data := map[string]string{
			"ps": tag,
			"port": m["port"],
			"id": m["id"],
			"aid": m["alterId"],
			"net": m["network"],
			"tls": m["security"],
			"add":m["add"],
			"v": "2",
		}

		if m["network"] == "http" {
			data["path"] = m["http.path"]
			data["net"] = "h2"
		} else if m["network"] == "ws" {
			data["path"] = m["ws.path"]
		} else if m["network"] == "kcp" {
			if m["kcp.type"] == "" {
				m["kcp.type"] = "none"
			}
			data["type"] = m["kcp.type"]
		} else if m["network"] == "quic" {
			if m["quic.type"] == "" {
				m["quic.type"] = "none"
			}
			data["type"] = m["quic.type"]
		}

		if data["security"] != "" {
			data["host"] = m["servername"]
		}

		strs, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			return "", err
		}

		return "vmess://" + base64.StdEncoding.EncodeToString(strs), nil
	},
}

var SSParser = JSONParser{
	Filed: map[string]FieldParser{
		"protocol":     JSONPathHandler("protocol"),
		"port":     JSONPathHandler("port"),
		"method":   JSONPathHandler("settings.method"),
		"password": JSONPathHandler("settings.password"),
	},
	DefaultField: map[string]string{
		"host": "",
	},
	PostHandler: func(m map[string]string, tag string) (string, error) {
		u := fmt.Sprintf("%s:%s@%s:%s", m["method"], m["password"], m["host"], m["port"])
		u = fmt.Sprintf("ss://%s#%s", base64.StdEncoding.EncodeToString([]byte(u)), tag)
		return u, nil
	},
}


