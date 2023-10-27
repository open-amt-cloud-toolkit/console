package companies

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

type CompaniesThing struct {
	router *http.ServeMux
	//parsed templates
	html *template.Template
}

func NewCompanies(router *http.ServeMux) CompaniesThing {
	//parse templates
	var err error
	html, err := templates.TemplateParseFSRecursive(templateFS, ".html", true, nil)
	if err != nil {
		panic(err)
	}

	ct := CompaniesThing{
		html: html,
	}
	router.Handle("/company/add", web.Action(ct.CompanyAdd))
	router.Handle("/company/add/", web.Action(ct.CompanyAdd))

	router.Handle("/company/edit", web.Action(ct.CompanyEdit))
	router.Handle("/company/edit/", web.Action(ct.CompanyEdit))

	router.Handle("/company", web.Action(ct.Companies))
	router.Handle("/company/", web.Action(ct.Companies))

	router.Handle("/companies", web.Action(ct.Index))

	return ct
}
func (ct CompaniesThing) Index(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, ct.html, "index.html", data, nil)
}

// GET /company/add
func (ct CompaniesThing) CompanyAdd(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, ct.html, "company-add.html", data, nil)
}

// /GET company/edit/{id}
func (ct CompaniesThing) CompanyEdit(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	row := ct.GetCompanyByID(id)
	return web.HTML(http.StatusOK, ct.html, "row-edit.html", row, nil)
}

// GET /company
// GET /company/{id}
// DELETE /company/{id}
// PUT /company/{id}
// POST /company
func (ct CompaniesThing) Companies(r *http.Request) *web.Response {
	id, segments := web.PathLast(r)
	switch r.Method {

	case http.MethodDelete:
		ct.DeleteCompany(id)
		return web.HTML(http.StatusOK, ct.html, "companies.html", data, nil)

	//cancel
	case http.MethodGet:
		if segments > 1 {
			//cancel edit
			row := ct.GetCompanyByID(id)
			return web.HTML(http.StatusOK, ct.html, "row.html", row, nil)
		} else {
			//cancel add
			return web.HTML(http.StatusOK, ct.html, "companies.html", data, nil)
		}

	//save edit
	case http.MethodPut:
		row := ct.GetCompanyByID(id)
		r.ParseForm()
		row.Company = r.Form.Get("company")
		row.Contact = r.Form.Get("contact")
		row.Country = r.Form.Get("country")
		ct.UpdateCompany(row)
		return web.HTML(http.StatusOK, ct.html, "row.html", row, nil)

	//save add
	case http.MethodPost:
		row := Company{}
		r.ParseForm()
		row.Company = r.Form.Get("company")
		row.Contact = r.Form.Get("contact")
		row.Country = r.Form.Get("country")
		ct.AddCompany(row)
		return web.HTML(http.StatusOK, ct.html, "companies.html", data, nil)
	}

	return web.Empty(http.StatusNotImplemented)
}
