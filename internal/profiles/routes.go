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
	return webtools.HTML(r, http.StatusOK, pt.html, "profiles/profiles-add.html", nil, nil)
}

func (pt ProfileThing) ProfileEdit(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	profile := pt.GetProfileByID(id)
	return webtools.HTML(r, http.StatusOK, pt.html, "profiles/profiles-edit.html", profile, nil)
}

func (pt ProfileThing) Profiles(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	switch r.Method {

	case http.MethodDelete:
		pt.DeleteProfile(id)
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/profiles.html", pt.GetProfiles(), nil)

	//cancel
	case http.MethodGet:
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/index.html", pt.GetProfiles(), nil)

	//save edit
	case http.MethodPut:
		profile := pt.GetProfileByID(id)
		r.ParseForm()
		profile.Id, _ = strconv.Atoi(id)
		profile.Name = r.Form.Get("name")
		profile.ControlMode = r.Form.Get("controlMode")
		if r.Form.Get("amtpassword") != "" {
			profile.AMTPassword = r.Form.Get("amtpassword")
		}
		if r.Form.Get("mebxpassword") != "" {
			profile.MEBXPassword = r.Form.Get("mebxpassword")
		}

		isValid, errors := profile.IsValid()
		if !isValid {
			return webtools.HTML(r, http.StatusBadRequest, pt.html, "profiles/errors.html", errors, nil)
		}
		pt.UpdateProfile(profile)
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/index.html", pt.GetProfiles(), nil)

	//save add
	case http.MethodPost:
		profile := Profile{}
		r.ParseForm()
		profile.Id, _ = strconv.Atoi(r.Form.Get("id"))
		profile.Name = r.Form.Get("name")
		profile.ControlMode = r.Form.Get("controlMode")
		profile.AMTPassword = r.Form.Get("amtpassword")
		profile.MEBXPassword = r.Form.Get("mebxpassword")

		isValid, errors := profile.IsValid()
		if !isValid {
			return webtools.HTML(r, http.StatusBadRequest, pt.html, "profiles/errors.html", errors, nil)
		}

		pt.AddDevice(profile)
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/index.html", pt.GetProfiles(), nil)
	}

	return web.Empty(http.StatusNotImplemented)
}
