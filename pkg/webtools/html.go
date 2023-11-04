package webtools

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jritsema/gotoolbox/web"
)

// Data returns a data response
func Data(status int, content []byte, headers web.Headers) *web.Response {
	return &web.Response{
		Status:  status,
		Content: bytes.NewBuffer(content),
		Headers: headers,
	}
}

// Empty returns an empty http response
func Empty(status int) *web.Response {
	return Data(status, []byte(""), nil)
}

// HTML renders an html template to a web response
func HTML(r *http.Request, status int, t *template.Template, template string, data interface{}, headers web.Headers) *web.Response {
	isHTMX := r.Header.Get("Hx-request")
	if isHTMX != "true" { // this is a full reload, need to return the full page
		template = "index.html"
	}
	//render template to buffer
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, template, data); err != nil {
		log.Println(err)
		return Empty(http.StatusInternalServerError)
	}
	return &web.Response{
		Status:      status,
		ContentType: "text/html",
		Content:     &buf,
		Headers:     headers,
	}
}

// ConvertPath takes in a path and returns it formatted with the appropriate path separator for the OS.
func ConvertPath(p string) string {
	// Split the string by either / or \ (so it works both ways)
	parts := strings.FieldsFunc(p, func(c rune) bool { return c == '/' || c == '\\' })
	return filepath.Join(parts...)
}
