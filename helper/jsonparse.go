package helper

import (
	"github.com/tidwall/gjson"
)

// FieldParser -
type FieldParser func(result gjson.Result) string

// JSONParser -
type JSONParser struct {
	TagFmt string
	Filed       map[string]FieldParser
	DefaultField map[string]string
	PostHandler func(m map[string]string, tag string) (string,error)
}

// Parse -
func (p *JSONParser) Parse(result gjson.Result) (string, error) {
	r := make(map[string]string)
	for k, v := range p.DefaultField {
		r[k]= v
	}

	for k, h := range p.Filed {
		r[k] = h(result)
	}

	tag := FmtStringReplace(p.TagFmt, r)
	return p.PostHandler(r, tag)
}

// JSONPathHandler -
func JSONPathHandler(p string) func(result gjson.Result) string {
	return func(result gjson.Result) string {
		return result.Get(p).String()
	}
}

