package dashboard

import (
	"html/template"
	"net/http"

	"github.com/jritsema/gotoolbox/web"
	"github.com/open-amt-cloud-toolkit/console/internal"
	"github.com/open-amt-cloud-toolkit/console/internal/i18n"
	"github.com/open-amt-cloud-toolkit/console/pkg/templates"
	"github.com/open-amt-cloud-toolkit/console/pkg/webtools"
)

type DashboardPages struct {
	//parsed templates
	html *template.Template
}

func NewDashboard(router *http.ServeMux) DashboardPages {
	funcMap := template.FuncMap{
		"Translate": i18n.Translate,
	}
	//parse templates
	var err error
	html, err := templates.TemplateParseFSRecursive(internal.TemplateFS, "/dashboard", ".html", true, funcMap)
	if err != nil {
		panic(err)
	}

	dp := DashboardPages{
		html: html,
	}

	router.Handle("/dashboard", web.Action(dp.Index))

	return dp
}

type DashboardContent struct{}

func (dp DashboardPages) Index(r *http.Request) *web.Response {
	dc := DashboardContent{}
	return webtools.HTML(r, http.StatusOK, dp.html, "dashboard/index.html", dc, nil)
}
