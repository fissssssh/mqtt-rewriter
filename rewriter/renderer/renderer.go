package renderer

import (
	"encoding/json"
	"strings"
	"text/template"
)

var renderers map[string]*template.Template = make(map[string]*template.Template)

func Render(name string, data interface{}, t string) (string, error) {
	if _, exist := renderers[name]; !exist {
		tt, err := template.New(name).Funcs(funcMap).Parse(t)
		if err != nil {
			return "", err
		}
		renderers[name] = tt
	}
	sb := strings.Builder{}
	err := renderers[name].Execute(&sb, data)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

var funcMap = template.FuncMap{"json": func(obj interface{}) string {
	jstring, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(jstring)
}}
