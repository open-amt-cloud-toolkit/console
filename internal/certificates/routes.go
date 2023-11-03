package certificates

import (
	"embed"
	"html/template"
	"net/http"

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
var (
	//go:embed all:templates/*
	templateFS embed.FS
)

type CertificateThing struct {
	router *http.ServeMux
	//parsed templates
	html *template.Template
}

func NewCertificates(router *http.ServeMux) CertificateThing {
	//parse templates
	var err error
	html, err := templates.TemplateParseFSRecursive(templateFS, ".html", true, nil)
	if err != nil {
		panic(err)
	}

	dt := CertificateThing{
		html: html,
	}
	router.Handle("/certificate/add", web.Action(dt.RouteAdd))
	router.Handle("/certificate/add/", web.Action(dt.RouteAdd))

	router.Handle("/certificate/edit", web.Action(dt.RouteEdit))
	router.Handle("/certificate/edit/", web.Action(dt.RouteEdit))

	router.Handle("/certificate", web.Action(dt.Certificates))
	router.Handle("/certificate/", web.Action(dt.Certificates))

	router.Handle("/certificates", web.Action(dt.Index))

	return dt
}
func (dt CertificateThing) Index(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, dt.html, "index.html", data, nil)
}

// GET /certificate/add
func (dt CertificateThing) RouteAdd(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, dt.html, "devices-add.html", data, nil)
}

// /GET company/edit/{id}
func (dt CertificateThing) RouteEdit(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	row := dt.GetByID(id)
	return web.HTML(http.StatusOK, dt.html, "row-edit.html", row, nil)
}

// GET /company
// GET /company/{id}
// DELETE /company/{id}
// PUT /company/{id}
// POST /company
func (dt CertificateThing) Certificates(r *http.Request) *web.Response {
	id, segments := web.PathLast(r)
	switch r.Method {

	case http.MethodDelete:
		dt.Delete(id)
		return web.HTML(http.StatusOK, dt.html, "devices.html", data, nil)

	//cancel
	case http.MethodGet:
		if segments > 1 {
			//cancel edit
			row := dt.GetByID(id)
			return web.HTML(http.StatusOK, dt.html, "row.html", row, nil)
		} else {
			//cancel add
			return web.HTML(http.StatusOK, dt.html, "devices.html", data, nil)
		}

	//save edit
	case http.MethodPut:
		row := dt.GetByID(id)
		r.ParseForm()
		row.UUID = id
		row.Name = r.Form.Get("name")
		row.IPAddress = r.Form.Get("ipaddress")
		row.FWVersion = r.Form.Get("fwversion")
		dt.Update(row)
		return web.HTML(http.StatusOK, dt.html, "row.html", row, nil)

	//save add
	case http.MethodPost:
		row := Certificate{}
		r.ParseForm()
		row.UUID = r.Form.Get("uuid")
		row.Name = r.Form.Get("name")
		row.IPAddress = r.Form.Get("ipaddress")
		row.FWVersion = r.Form.Get("fwversion")
		dt.Add(row)
		return web.HTML(http.StatusOK, dt.html, "devices.html", data, nil)
	}

	return web.Empty(http.StatusNotImplemented)
}
