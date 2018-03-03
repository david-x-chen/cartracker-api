package common

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// StringInSlice is the utility
func StringInSlice(a string, list []string) bool {
	var result bool
	for _, b := range list {
		result = strings.EqualFold(b, a)
		if result {
			break
		}
	}
	return result
}

// RenderTemplate Render a template, or server error.
func RenderTemplate(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
		fmt.Printf("\nRender Error: %v\n", err)
		return
	}
	w.Write(buf.Bytes())
}

// PushStaticResource Push the given resource to the client.
func PushStaticResource(w http.ResponseWriter, resource string) {
	pusher, ok := w.(http.Pusher)
	if ok {
		if err := pusher.Push(resource, nil); err == nil {
			return
		}
	}
}
