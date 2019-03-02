package helper

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

const (
	defaultVmessName = "{protocol}-{network}"
	defaultSSName    = "{protocol}"
)

// from https://github.com/2dust/v2rayN/wiki/分享链接格式说明(ver-2)
var vmessParser = JSONParser{
	Filed: map[string]FieldParser{
		"protocol":     JSONPathHandler("protocol"),
		"port": JSONPathHandler("port"),
		"id":   JSONPathHandler("settings.clients.0.id"),
		"alterId":  JSONPathHandler("settings.clients.0.alterId"),
		"network":  JSONPathHandler("streamSettings.network"),
		"security":  JSONPathHandler("streamSettings.security"),
		"http.host": JSONPathHandler("streamSettings.httpSettings.host.0"),
		"http.path": JSONPathHandler("streamSettings.httpSettings.path"),
		"ws.path": JSONPathHandler("streamSettings.wsSettings.path"),
	},
	DefaultField: map[string]string{
		"v":   "2",
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
		}

		if m["network"] == "http" {
			data["host"] = m["http.host"]
			data["path"] = m["http.path"]
		} else if m["network"] == "ws" {
			data["path"] = m["ws.path"]
		}

		strs, err := json.Marshal(data)
		if err != nil {
			return "", err
		}

		return "vmess://" + base64.StdEncoding.EncodeToString(strs), nil
	},
}

var ssParser = JSONParser{
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

// Export -
func Export(filepath, host, vmessfmt, ssfmt string) ([]string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err !=nil {
		return nil, err
	}

	value := gjson.Get(string(data), "inbounds")
	if !value.IsArray() {
		return nil, errors.New("Unknow Format")
	}

	ret := make([]string, 0)

	for _, inbound := range value.Array() {

		protocol := inbound.Get("protocol").String()
		if protocol == "vmess" {
			vmessParser.DefaultField["add"] = host
			vmessParser.TagFmt = vmessfmt
			u, err := vmessParser.Parse(inbound)
			if err != nil {
				return nil, err
			}
			ret = append(ret, u)
		} else if protocol == "shadowsocks" {
			ssParser.DefaultField["host"] = host
			ssParser.TagFmt = ssfmt
			u, err := ssParser.Parse(inbound)
			if err != nil {
				return nil, err
			}
			ret = append(ret, u)
		}
	}


	return ret, nil
}
