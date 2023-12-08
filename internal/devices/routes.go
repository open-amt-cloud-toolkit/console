package devices

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/jritsema/go-htmx-starter/internal"
	"github.com/jritsema/go-htmx-starter/pkg/templates"
	"github.com/jritsema/go-htmx-starter/pkg/webtools"
	"github.com/jritsema/gotoolbox/web"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/power"
	"go.etcd.io/bbolt"
)

type DeviceThing struct {
	db *bbolt.DB
	//parsed templates
	html *template.Template
}

func NewDevices(db *bbolt.DB, router *http.ServeMux) DeviceThing {
	//parse templates
	var err error

	funcMap := template.FuncMap{
		"ProvisioningModeLookup":  ProvisioningModeLookup,
		"ProvisioningStateLookup": ProvisioningStateLookup,
	}
	html, err := templates.TemplateParseFSRecursive(internal.TemplateFS, "/devices", ".html", true, funcMap)
	if err != nil {
		panic(err)
	}

	// html.Funcs(funcMap)

	dt := DeviceThing{
		db:   db,
		html: html,
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Devices"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
	router.Handle("/device/add", web.Action(dt.DeviceAdd))
	router.Handle("/device/add/", web.Action(dt.DeviceAdd))

	router.Handle("/device/edit", web.Action(dt.DeviceEdit))
	router.Handle("/device/edit/", web.Action(dt.DeviceEdit))

	router.Handle("/device", web.Action(dt.Devices))
	router.Handle("/device/", web.Action(dt.Devices))

	router.Handle("/devices", web.Action(dt.Index))

	router.Handle("/device/connect", web.Action(dt.DeviceConnect))
	router.Handle("/device/connect/", web.Action(dt.DeviceConnect))

	router.Handle("/device/ethernet", web.Action(dt.GetEthernet))
	router.Handle("/device/ethernet/", web.Action(dt.GetEthernet))

	router.Handle("/device/powerState/", web.Action(dt.ChangePowerState))

	router.Handle("/device/wsman-explorer", web.Action(dt.GetWsmanExplorer))
	router.Handle("/device/wsman-explorer/", web.Action(dt.GetWsmanExplorer))

	router.Handle("/device/ws-classes", web.Action(dt.GetWsmanClasses))
	router.Handle("/device/ws-classes/", web.Action(dt.GetWsmanClasses))

	router.Handle("/device/ws-methods", web.Action(dt.GetWsmanMethods))
	router.Handle("/device/ws-methods/", web.Action(dt.GetWsmanMethods))

	router.Handle("/device/test", web.Action(dt.WsmanTest))
	router.Handle("/device/test/", web.Action(dt.WsmanTest))

	return dt
}
func (dt DeviceThing) Index(r *http.Request) *web.Response {
	return webtools.HTML(r, http.StatusOK, dt.html, "devices/index.html", dt.GetDevices(), nil)
}

// GET /device/add
func (dt DeviceThing) DeviceAdd(r *http.Request) *web.Response {
	return webtools.HTML(r, http.StatusOK, dt.html, "devices/devices-add.html", dt.GetDevices(), nil)
}

// /GET device/edit/{id}
func (dt DeviceThing) DeviceEdit(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	row := dt.GetDeviceByID(id)
	return webtools.HTML(r, http.StatusOK, dt.html, "devices/row-edit.html", row, nil)
}

func (dt DeviceThing) GetWsmanExplorer(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	device := dt.GetDeviceByID(id)
	return webtools.HTML(r, http.StatusOK, dt.html, "devices/wsman-explorer/wsman-explorer.html", device, nil)
}

func (dt DeviceThing) GetWsmanClasses(r *http.Request) *web.Response {
	classes := GetSupportedWsmanClasses("")
	return webtools.HTML(r, http.StatusOK, dt.html, "devices/wsman-explorer/class-select.html", classes, nil)
}

func (dt DeviceThing) GetWsmanMethods(r *http.Request) *web.Response {
	queryValues := r.URL.Query()
	selected := queryValues.Get("class-selector")
	class := GetSupportedWsmanClasses(selected)
	methods := class[0].MethodList
	return webtools.HTML(r, http.StatusOK, dt.html, "devices/wsman-explorer/method-select.html", methods, nil)
}

func (dt DeviceThing) WsmanTest(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	device := dt.GetDeviceByID(id)
	r.ParseForm()
	class := r.Form.Get("class-selector")
	method := r.Form.Get("method-selector")

	response, err := MakeWsmanCall(device, class, method)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return webtools.HTML(r, http.StatusOK, dt.html, "devices/wsman-explorer/wsman.html", response, nil)
}

// Connect to device
func (dt DeviceThing) DeviceConnect(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	device := dt.GetDeviceByID(id)
	wsman := CreateWsmanConnection(device)
	// Get General Settings
	gs, err := GetGeneralSettings(wsman)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Get Setup and Configuration Service
	scs, err := GetSetupAndConfigurationService(wsman)
	if err != nil {
		fmt.Println("Error:", err)
	}

	uuid, err := GetDeviceUUID(wsman)
	if err != nil {
		fmt.Println("Error:", err)
	}

	device.AMTSpecific.UUID = uuid

	dc := DeviceContent{
		Device:                       device,
		GeneralSettings:              gs,
		SetupAndConfigurationService: scs,
	}

	return webtools.HTML(r, http.StatusOK, dt.html, "devices/device.html", dc, nil)
}

func (dt DeviceThing) GetEthernet(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	queryValues := r.URL.Query()
	keyValue := queryValues.Get("eth")
	eth, err := strconv.Atoi(keyValue)
	if err != nil {
		fmt.Println("Error:", err)
	}
	device := dt.GetDeviceByID(id)
	wsman := CreateWsmanConnection(device)
	// Get Ethernet Settings
	ep, err := GetEthernetSettings(wsman, eth)
	if err != nil {
		fmt.Println("Error:", err)
	}
	ec := EthernetContent{
		EthernetPort: ep,
	}
	if ec.EthernetPort.ElementName == "" {
		return webtools.HTML(r, http.StatusOK, dt.html, "devices/ethernet.html", nil, nil)
	}
	return webtools.HTML(r, http.StatusOK, dt.html, "devices/ethernet.html", ec.EthernetPort, nil)
}

func (dt DeviceThing) ChangePowerState(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	queryValues := r.URL.Query()
	keyValue := queryValues.Get("power")
	technology := "amt"
	powerStateRequested := getPowerStateValue(technology, keyValue)
	device := dt.GetDeviceByID(id)
	wsman := CreateWsmanConnection(device)
	response, errors := ChangePowerState(wsman, power.PowerState(powerStateRequested))
	if errors != nil {
		return webtools.HTML(r, http.StatusRequestTimeout, dt.html, "devices/errors.html", errors, nil)
	}
	return webtools.HTML(r, http.StatusOK, dt.html, "devices/device.html", response, nil)
}

// GET /device
// GET /device/{id}
// DELETE /device/{id}
// PUT /device/{id}
// POST /device
func (dt DeviceThing) Devices(r *http.Request) *web.Response {
	id, segments := web.PathLast(r)
	switch r.Method {

	case http.MethodDelete:
		dt.DeleteDevice(id)
		return webtools.HTML(r, http.StatusOK, dt.html, "devices/devices.html", dt.GetDevices(), nil)

	//cancel
	case http.MethodGet:
		if segments > 1 {
			//cancel edit
			row := dt.GetDeviceByID(id)
			return webtools.HTML(r, http.StatusOK, dt.html, "devices/row.html", row, nil)
		} else {
			//cancel add
			return webtools.HTML(r, http.StatusOK, dt.html, "devices/devices.html", dt.GetDevices(), nil)
		}

	//save edit
	case http.MethodPut:
		row := dt.GetDeviceByID(id)
		r.ParseForm()
		row.Id, _ = strconv.Atoi(id)
		row.Name = r.Form.Get("name")
		row.Address = r.Form.Get("address")
		row.Username = r.Form.Get("username")
		row.Password = r.Form.Get("password")
		tls := false
		if r.Form.Get("usetls") == "on" {
			tls = true
		}
		row.UseTLS = tls
		selfSignedAllowed := false
		if r.Form.Get("selfsignedallowed") == "on" {
			selfSignedAllowed = true
		}
		row.SelfSignedAllowed = selfSignedAllowed
		isValid, errors := row.IsValid()
		if !isValid {
			return webtools.HTML(r, http.StatusBadRequest, dt.html, "devices/errors.html", errors, nil)
		}
		dt.UpdateDevice(row)
		return webtools.HTML(r, http.StatusOK, dt.html, "devices/row.html", row, nil)

	//save add
	case http.MethodPost:
		row := Device{}
		r.ParseForm()
		row.Id, _ = strconv.Atoi(r.Form.Get("id"))
		row.Name = r.Form.Get("name")
		row.Address = r.Form.Get("address")
		row.Username = r.Form.Get("username")
		row.Password = r.Form.Get("password")
		tls := false
		if r.Form.Get("usetls") == "on" {
			tls = true
		}
		row.UseTLS = tls
		selfSignedAllowed := false
		if r.Form.Get("selfsignedallowed") == "on" {
			selfSignedAllowed = true
		}
		row.SelfSignedAllowed = selfSignedAllowed

		isValid, errors := row.IsValid()
		if !isValid {
			return webtools.HTML(r, http.StatusBadRequest, dt.html, "devices/errors.html", errors, nil)
		}

		dt.AddDevice(row)
		return webtools.HTML(r, http.StatusOK, dt.html, "devices/devices.html", dt.GetDevices(), nil)
	}

	return web.Empty(http.StatusNotImplemented)
}
