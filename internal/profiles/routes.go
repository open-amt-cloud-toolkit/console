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
	router.Handle("/profile/export", web.Action(pt.ExportProfile))
	router.Handle("/profile/export/", web.Action(pt.ExportProfile))
	router.Handle("/profile/download", web.Action(pt.Download))
	router.Handle("/profile/download/", web.Action(pt.Download))
	router.Handle("/profile/technology-select/edit", web.Action(pt.TechnologySelectEdit))
	router.Handle("/profile/technology-select/edit/", web.Action(pt.TechnologySelectEdit))
	router.Handle("/profile/technology-select/add", web.Action(pt.TechnologySelectAdd))
	router.Handle("/profile/technology-select/add/", web.Action(pt.TechnologySelectAdd))
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

const (
	AMT     = "AMT"
	BMC     = "BMC"
	DASH    = "DASH"
	Redfish = "Redfish"
)

func (pt ProfileThing) TechnologySelectAdd(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	profile := pt.GetProfileByID(id)
	queryValues := r.URL.Query()
	technologyValue := queryValues.Get("technology")
	switch technologyValue {
	case AMT:
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/amt/add.html", profile, nil)
	case BMC:
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/bmc/add.html", profile, nil)
	case DASH:
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/dash/add.html", profile, nil)
	case Redfish:
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/redfish/add.html", profile, nil)
	default:
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/no-technology.html", nil, nil)
	}
}

func (pt ProfileThing) TechnologySelectEdit(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	profile := pt.GetProfileByID(id)
	queryValues := r.URL.Query()
	technologyValue := queryValues.Get("technology")
	switch technologyValue {
	case AMT:
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/amt/edit.html", profile, nil)
	case BMC:
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/bmc/edit.html", profile, nil)
	case DASH:
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/dash/edit.html", profile, nil)
	case Redfish:
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/redfish/edit.html", profile, nil)
	default:
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/no-technology.html", nil, nil)
	}
}

func (pt ProfileThing) Profiles(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	switch r.Method {

	case http.MethodDelete:
		pt.DeleteProfile(id)
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/profiles.html", pt.GetProfiles(), nil)

	//Gets
	case http.MethodGet:
		// Cancel
		if id == "profiles" {
			return webtools.HTML(r, http.StatusOK, pt.html, "profiles/index.html", pt.GetProfiles(), nil)
		}

	//save edit
	case http.MethodPut:
		profile := pt.GetProfileByID(id)
		r.ParseForm()
		profile.Id, _ = strconv.Atoi(id)
		profile.Name = r.Form.Get("name")
		profile.Configuration.AMTSpecific.ControlMode = r.Form.Get("controlMode")
		if r.Form.Get("adminPassword") != "" {
			profile.Configuration.RemoteManagement.AdminPassword = r.Form.Get("adminPassword")
		}
		if r.Form.Get("mebxPassword") != "" {
			profile.Configuration.AMTSpecific.MEBXPassword = r.Form.Get("mebxPassword")
		}
		profile.Configuration.RemoteManagement.GeneralSettings.HostName = r.Form.Get("hostName")
		profile.Configuration.RemoteManagement.GeneralSettings.DomainName = r.Form.Get("domainName")
		profile.Configuration.RemoteManagement.GeneralSettings.NetworkEnabled = checkboxValue(r.Form.Get("networkEnabled"))
		profile.Configuration.RemoteManagement.GeneralSettings.SharedFQDN = checkboxValue(r.Form.Get("sharedFQDN"))
		profile.Configuration.RemoteManagement.GeneralSettings.PingResponseEnabled = checkboxValue(r.Form.Get("pingResponseEnabled"))

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
		profile.Technology = r.Form.Get("technology")
		profile.Configuration.AMTSpecific.ControlMode = r.Form.Get("controlMode")
		profile.Configuration.RemoteManagement.AdminPassword = r.Form.Get("adminPassword")
		profile.Configuration.AMTSpecific.MEBXPassword = r.Form.Get("mebxPassword")
		profile.Configuration.RemoteManagement.GeneralSettings.HostName = r.Form.Get("hostName")
		profile.Configuration.RemoteManagement.GeneralSettings.DomainName = r.Form.Get("domainName")
		profile.Configuration.RemoteManagement.GeneralSettings.NetworkEnabled = checkboxValue(r.Form.Get("networkEnabled"))
		profile.Configuration.RemoteManagement.GeneralSettings.SharedFQDN = checkboxValue(r.Form.Get("sharedFQDN"))
		profile.Configuration.RemoteManagement.GeneralSettings.PingResponseEnabled = checkboxValue(r.Form.Get("pingResponseEnabled"))

		isValid, errors := profile.IsValid()
		if !isValid {
			return webtools.HTML(r, http.StatusBadRequest, pt.html, "profiles/errors.html", errors, nil)
		}

		err := pt.AddProfile(profile)
		if err != nil {
			fmt.Println(err)
			return webtools.HTML(r, http.StatusInternalServerError, pt.html, "profiles/errors.html", err, nil)
		}
		return webtools.HTML(r, http.StatusOK, pt.html, "profiles/index.html", pt.GetProfiles(), nil)
	}

	return web.Empty(http.StatusNotImplemented)
}
