package dashboard

import (
	"html/template"
	"net/http"

	"github.com/jritsema/go-htmx-starter/internal"
	"github.com/jritsema/go-htmx-starter/pkg/templates"
	"github.com/jritsema/go-htmx-starter/pkg/webtools"
	"github.com/jritsema/gotoolbox/web"
)

type DashboardPages struct {
	router *http.ServeMux
	//parsed templates
	html *template.Template
}

func NewDashboard(router *http.ServeMux) DashboardPages {

	//parse templates
	var err error
	html, err := templates.TemplateParseFSRecursive(internal.TemplateFS, "/dashboard", ".html", true, nil)
	if err != nil {
		panic(err)
	}

	dt := DashboardPages{
		html: html,
	}
	// router.Handle("/profile/add", web.Action(dt.RouteAdd))
	// router.Handle("/profile/add/", web.Action(dt.RouteAdd))

	// router.Handle("/profile/edit", web.Action(dt.RouteEdit))
	// router.Handle("/profile/edit/", web.Action(dt.RouteEdit))

	// router.Handle("/profile", web.Action(dt.Profiles))
	// router.Handle("/profile/", web.Action(dt.Profiles))

	router.Handle("/dashboard", web.Action(dt.Index))

	return dt
}
func (dt DashboardPages) Index(r *http.Request) *web.Response {
	return webtools.HTML(r, http.StatusOK, dt.html, "dashboard/index.html", nil, nil)
}
