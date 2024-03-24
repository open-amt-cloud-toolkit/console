package certificates

import (
	"html/template"
	"net/http"

	"github.com/jritsema/gotoolbox/web"
	"github.com/open-amt-cloud-toolkit/console/internal"
	"github.com/open-amt-cloud-toolkit/console/internal/i18n"
	"github.com/open-amt-cloud-toolkit/console/pkg/templates"
	"github.com/open-amt-cloud-toolkit/console/pkg/webtools"
)

// Delete -> DELETE /company/{id} -> delete, companys.html

// Edit   -> GET /company/edit/{id} -> row-edit.html
// Save   ->   PUT /company/{id} -> update, row.html
// Cancel ->	 GET /company/{id} -> nothing, row.html

// Add    -> GET /company/add/ -> companys-add.html (target body with row-add.html and row.html)
// Save   ->   POST /company -> add, companys.html (target body without row-add.html)
// Cancel ->	 GET /company -> nothing, companys.html

type CertificateThing struct {
	//parsed templates
	html *template.Template
}

func NewCertificates(router *http.ServeMux) CertificateThing {
	funcMap := template.FuncMap{
		"Translate": i18n.Translate,
	}
	//parse templates
	var err error
	html, err := templates.TemplateParseFSRecursive(internal.TemplateFS, "/certificates", ".html", true, funcMap)
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
	return webtools.HTML(r, http.StatusOK, dt.html, "certificates/index.html", data, nil)
}

// GET /certificate/add
func (dt CertificateThing) RouteAdd(r *http.Request) *web.Response {
	return webtools.HTML(r, http.StatusOK, dt.html, "devices-add.html", data, nil)
}

// /GET company/edit/{id}
func (dt CertificateThing) RouteEdit(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	row := dt.GetByID(id)
	return webtools.HTML(r, http.StatusOK, dt.html, "row-edit.html", row, nil)
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
		return webtools.HTML(r, http.StatusOK, dt.html, "devices.html", data, nil)

	//cancel
	case http.MethodGet:
		if segments > 1 {
			//cancel edit
			row := dt.GetByID(id)
			return webtools.HTML(r, http.StatusOK, dt.html, "row.html", row, nil)
		} else {
			//cancel add
			return webtools.HTML(r, http.StatusOK, dt.html, "devices.html", data, nil)
		}

	//save edit
	case http.MethodPut:
		row := dt.GetByID(id)
		_ = r.ParseForm()
		row.UUID = id
		row.Name = r.Form.Get("name")
		row.IPAddress = r.Form.Get("ipaddress")
		row.FWVersion = r.Form.Get("fwversion")
		dt.Update(row)
		return webtools.HTML(r, http.StatusOK, dt.html, "row.html", row, nil)

	//save add
	case http.MethodPost:
		row := Certificate{}
		_ = r.ParseForm()
		row.UUID = r.Form.Get("uuid")
		row.Name = r.Form.Get("name")
		row.IPAddress = r.Form.Get("ipaddress")
		row.FWVersion = r.Form.Get("fwversion")
		dt.Add(row)
		return webtools.HTML(r, http.StatusOK, dt.html, "devices.html", data, nil)
	}

	return web.Empty(http.StatusNotImplemented)
}
