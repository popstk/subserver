package helper

import (
	"github.com/tidwall/gjson"
)

// FieldParser -
type FieldParser func(result gjson.Result) string

// JSONParser -
type JSONParser struct {
	TagFmt       string
	Filed        map[string]FieldParser
	DefaultField map[string]string
	PostHandler  func(m map[string]string, tag string) (Endpoint, error)
}

// Parse -
func (p *JSONParser) Parse(result gjson.Result) (Endpoint, error) {
	r := make(map[string]string)
	for k, h := range p.Filed {
		r[k] = h(result)
	}

	for k, v := range p.DefaultField {
		r[k] = v
	}

	tag := FmtStringReplace(p.TagFmt, r)
	return p.PostHandler(r, tag)
}

// JSONPath -
func JSONPath(p string) func(result gjson.Result) string {
	return func(result gjson.Result) string {
		return result.Get(p).String()
	}
}

// JSONPathWithValue json path with default value
func JSONPathWithValue(p, value string) func(result gjson.Result) string {
	return func(result gjson.Result) string {
		if v := result.Get(p).String(); v != "" {
			return v
		}
		return value
	}
}
