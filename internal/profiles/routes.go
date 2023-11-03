package profiles

import (
	"html/template"
	"net/http"

	"github.com/jritsema/go-htmx-starter/internal"
	"github.com/jritsema/go-htmx-starter/pkg/templates"
	"github.com/jritsema/gotoolbox/web"
)

// Delete -> DELETE /company/{id} -> delete, companys.html

// Edit   -> GET /company/edit/{id} -> row-edit.html
// Save   ->   PUT /company/{id} -> update, row.html
// Cancel ->	 GET /company/{id} -> nothing, row.html

// Add    -> GET /company/add/ -> companys-add.html (target body with row-add.html and row.html)
// Save   ->   POST /company -> add, companys.html (target body without row-add.html)
// Cancel ->	 GET /company -> nothing, companys.html

type CertificateThing struct {
	router *http.ServeMux
	//parsed templates
	html *template.Template
}

func NewProfiles(router *http.ServeMux) CertificateThing {
	//parse templates
	var err error
	html, err := templates.TemplateParseFSRecursive(internal.TemplateFS, ".html", true, nil)
	if err != nil {
		panic(err)
	}

	dt := CertificateThing{
		html: html,
	}
	// router.Handle("/profile/add", web.Action(dt.RouteAdd))
	// router.Handle("/profile/add/", web.Action(dt.RouteAdd))

	// router.Handle("/profile/edit", web.Action(dt.RouteEdit))
	// router.Handle("/profile/edit/", web.Action(dt.RouteEdit))

	// router.Handle("/profile", web.Action(dt.Profiles))
	// router.Handle("/profile/", web.Action(dt.Profiles))

	router.Handle("/profiles", web.Action(dt.Index))

	return dt
}
func (dt CertificateThing) Index(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, dt.html, "index.html", nil, nil)
}
