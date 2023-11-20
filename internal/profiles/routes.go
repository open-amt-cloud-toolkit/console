package profiles

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/jritsema/go-htmx-starter/internal"
	"github.com/jritsema/go-htmx-starter/pkg/templates"
	"github.com/jritsema/go-htmx-starter/pkg/webtools"
	"github.com/jritsema/gotoolbox/web"
	"go.etcd.io/bbolt"
)

// Delete -> DELETE /company/{id} -> delete, companys.html

// Edit   -> GET /company/edit/{id} -> row-edit.html
// Save   ->   PUT /company/{id} -> update, row.html
// Cancel ->	 GET /company/{id} -> nothing, row.html

// Add    -> GET /company/add/ -> companys-add.html (target body with row-add.html and row.html)
// Save   ->   POST /company -> add, companys.html (target body without row-add.html)
// Cancel ->	 GET /company -> nothing, companys.html

type ProfileThing struct {
	db *bbolt.DB
	//parsed templates
	html *template.Template
}

func NewProfiles(db *bbolt.DB, router *http.ServeMux) ProfileThing {

	//parse templates
	var err error
	html, err := templates.TemplateParseFSRecursive(internal.TemplateFS, "/profiles", ".html", true, nil)
	if err != nil {
		panic(err)
	}

	pt := ProfileThing{
		db:   db,
		html: html,
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Profiles"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	router.Handle("/profile/add", web.Action(pt.ProfileAdd))
	router.Handle("/profile/add/", web.Action(pt.ProfileAdd))

	router.Handle("/profile/edit", web.Action(pt.ProfileEdit))
	router.Handle("/profile/edit/", web.Action(pt.ProfileEdit))

	router.Handle("/profile", web.Action(pt.Profiles))
	router.Handle("/profile/", web.Action(pt.Profiles))

	router.Handle("/profiles", web.Action(pt.Index))

	return pt
}
func (pt ProfileThing) Index(r *http.Request) *web.Response {
	return webtools.HTML(r, http.StatusOK, pt.html, "profiles/index.html", pt.GetProfiles(), nil)
}

func (pt ProfileThing) ProfileAdd(r *http.Request) *web.Response {
	return webtools.HTML(r, http.StatusOK, pt.html, "profiles/profiles-add.html", pt.GetProfiles(), nil)
}

func (pt ProfileThing) ProfileEdit(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	row := pt.GetProfileByID(id)
	return webtools.HTML(r, http.StatusOK, pt.html, "profiles/row-edit.html", row, nil)
}

func (pt ProfileThing) Profiles(r *http.Request) *web.Response {
	id, segments := web.PathLast(r)
	switch r.Method {

	case http.MethodDelete:
		pt.DeleteProfile(id)
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/profiles.html", pt.GetProfiles(), nil)

	//cancel
	case http.MethodGet:
		if segments > 1 {
			//cancel edit
			row := pt.GetProfileByID(id)
			return webtools.HTML(r, http.StatusOK, pt.html, "profiles/row.html", row, nil)
		} else {
			//cancel add
			return webtools.HTML(r, http.StatusOK, pt.html, "profiles/profiles.html", pt.GetProfiles(), nil)
		}

	//save edit
	case http.MethodPut:
		row := pt.GetProfileByID(id)
		r.ParseForm()
		row.Id, _ = strconv.Atoi(id)
		row.Name = r.Form.Get("name")
		isValid, errors := row.IsValid()
		if !isValid {
			return webtools.HTML(r, http.StatusBadRequest, pt.html, "profiles/errors.html", errors, nil)
		}
		pt.UpdateProfile(row)
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/row.html", row, nil)

	//save add
	case http.MethodPost:
		row := Profile{}
		r.ParseForm()
		row.Id, _ = strconv.Atoi(r.Form.Get("id"))
		row.Name = r.Form.Get("name")

		isValid, errors := row.IsValid()
		if !isValid {
			return webtools.HTML(r, http.StatusBadRequest, pt.html, "profiles/errors.html", errors, nil)
		}

		pt.AddDevice(row)
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/profiles.html", pt.GetProfiles(), nil)
	}

	return web.Empty(http.StatusNotImplemented)
}
