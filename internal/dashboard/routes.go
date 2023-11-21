package dashboard

import (
	"html/template"
	"net/http"

	"github.com/jritsema/go-htmx-starter/internal"
	"github.com/jritsema/go-htmx-starter/internal/devices"
	"github.com/jritsema/go-htmx-starter/internal/profiles"
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

	dp := DashboardPages{
		html: html,
	}

	router.Handle("/dashboard", web.Action(dp.Index))

	return dp
}

type DashboardContent struct {
	devices []devices.Device
	profiles []profiles.Profile
}

func (dp DashboardPages) Index(r *http.Request) *web.Response {
	dc := DashboardContent{}
	dt := devices.DeviceThing{}
	dc.devices = dt.GetDevices()
	pt := profiles.ProfileThing{}
	dc.profiles = pt.GetProfiles()
	return webtools.HTML(r, http.StatusOK, dp.html, "dashboard/index.html", dc, nil)
}
