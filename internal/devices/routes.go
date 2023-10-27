package devices

import (
	"embed"
	"html/template"
	"net/http"

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

type DeviceThing struct {
	router *http.ServeMux
	//parsed templates
	html *template.Template
}

func NewDevices(router *http.ServeMux) DeviceThing {
	//parse templates
	var err error
	html, err := web.TemplateParseFSRecursive(templateFS, ".html", true, nil)
	if err != nil {
		panic(err)
	}

	dt := DeviceThing{
		html: html,
	}
	router.Handle("/device/add", web.Action(dt.DeviceAdd))
	router.Handle("/device/add/", web.Action(dt.DeviceAdd))

	router.Handle("/device/edit", web.Action(dt.DeviceEdit))
	router.Handle("/device/edit/", web.Action(dt.DeviceEdit))

	router.Handle("/device", web.Action(dt.Devices))
	router.Handle("/device/", web.Action(dt.Devices))

	router.Handle("/devices", web.Action(dt.Index))

	return dt
}
func (dt DeviceThing) Index(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, dt.html, "index.html", data, nil)
}

// GET /device/add
func (dt DeviceThing) DeviceAdd(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, dt.html, "devices-add.html", data, nil)
}

// /GET company/edit/{id}
func (dt DeviceThing) DeviceEdit(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	row := dt.GetDeviceByID(id)
	return web.HTML(http.StatusOK, dt.html, "row-edit.html", row, nil)
}

// GET /company
// GET /company/{id}
// DELETE /company/{id}
// PUT /company/{id}
// POST /company
func (dt DeviceThing) Devices(r *http.Request) *web.Response {
	id, segments := web.PathLast(r)
	switch r.Method {

	case http.MethodDelete:
		dt.DeleteDevice(id)
		return web.HTML(http.StatusOK, dt.html, "devices.html", data, nil)

	//cancel
	case http.MethodGet:
		if segments > 1 {
			//cancel edit
			row := dt.GetDeviceByID(id)
			return web.HTML(http.StatusOK, dt.html, "row.html", row, nil)
		} else {
			//cancel add
			return web.HTML(http.StatusOK, dt.html, "devices.html", data, nil)
		}

	//save edit
	case http.MethodPut:
		row := dt.GetDeviceByID(id)
		r.ParseForm()
		row.UUID = id
		row.Name = r.Form.Get("name")
		row.IPAddress = r.Form.Get("ipaddress")
		row.FWVersion = r.Form.Get("fwversion")
		dt.UpdateDevice(row)
		return web.HTML(http.StatusOK, dt.html, "row.html", row, nil)

	//save add
	case http.MethodPost:
		row := Device{}
		r.ParseForm()
		row.UUID = r.Form.Get("uuid")
		row.Name = r.Form.Get("name")
		row.IPAddress = r.Form.Get("ipaddress")
		row.FWVersion = r.Form.Get("fwversion")
		dt.AddDevice(row)
		return web.HTML(http.StatusOK, dt.html, "devices.html", data, nil)
	}

	return web.Empty(http.StatusNotImplemented)
}
