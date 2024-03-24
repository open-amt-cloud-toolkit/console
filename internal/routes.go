package internal

import (
	"embed"
	"html/template"
	"net/http"
	"os"

	"github.com/jritsema/gotoolbox/web"
	"github.com/open-amt-cloud-toolkit/console/internal/i18n"
	"github.com/open-amt-cloud-toolkit/console/pkg/templates"
	"github.com/open-amt-cloud-toolkit/console/pkg/webtools"
)

type IndexThing struct {
}

var (
	//parsed templates
	html *template.Template
	//go:embed all:templates/**
	TemplateFS embed.FS
)

func NewIndex(router *http.ServeMux) IndexThing {

	funcMap := template.FuncMap{
		"Translate": i18n.Translate,
	}

	//parse templates
	var err error
	html, err = templates.TemplateParseFSRecursive(TemplateFS, "/", ".html", true, funcMap)
	if err != nil {
		panic(err)
	}

	it := IndexThing{}
	router.Handle("/", web.Action(it.Index))
	router.Handle("/index.html", web.Action(it.Index))
	router.Handle("/menu", web.Action(it.Menu))
	router.Handle("/close", web.Action(it.Close))

	return it
}

func (it IndexThing) Index(r *http.Request) *web.Response {
	return webtools.HTML(r, http.StatusOK, html, "index.html", nil, nil)
}

func (it IndexThing) Menu(r *http.Request) *web.Response {
	return webtools.HTML(r, http.StatusOK, html, "menu.html", nil, nil)
}

func (it IndexThing) Close(r *http.Request) *web.Response {
	os.Exit(0)
	return webtools.HTML(r, http.StatusOK, html, "", nil, nil)
}
