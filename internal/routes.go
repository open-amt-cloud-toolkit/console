package internal

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/jritsema/go-htmx-starter/pkg/templates"
	"github.com/jritsema/gotoolbox/web"
)

type IndexThing struct {
}

var (
	router *http.ServeMux
	//parsed templates
	html *template.Template
	//go:embed all:templates/**
	TemplateFS embed.FS
)

func NewIndex(router *http.ServeMux) IndexThing {
	//parse templates
	var err error
	html, err = templates.TemplateParseFSRecursive(TemplateFS, ".html", true, nil)
	if err != nil {
		panic(err)
	}

	it := IndexThing{}
	router.Handle("/", web.Action(it.Index))
	router.Handle("/index.html", web.Action(it.Index))
	router.Handle("/menu", web.Action(it.Menu))

	return it
}

func (it IndexThing) Index(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, html, "index.html", nil, nil)
}

func (it IndexThing) Menu(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, html, "menu.html", nil, nil)
}
