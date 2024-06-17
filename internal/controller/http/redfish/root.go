package redfish

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/redfish"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type deviceManagementRoutes struct {
	d devices.Feature
	l logger.Interface
}

func NewRedfishRoutes(handler *gin.RouterGroup, d devices.Feature, l logger.Interface) {
	r := &deviceManagementRoutes{d, l}

	h := handler.Group("")
	{
		h.GET("$metadata", r.getMetadata)
		h.GET("", r.getServiceRoot)
		h.GET("SessionService", r.getSessionService)
		h.GET("SessionService/Sessions", r.getSessionCollection)
		h.GET("SessionService/Sessions/:id", r.getSession)
		h.GET("Systems/:id", r.getComputerSystem)
		h.GET("Systems", r.getComputerSystemCollection)

	}
}

// Initialize ServiceRoot
func (r *deviceManagementRoutes) getServiceRoot(c *gin.Context) {
	serviceRoot := redfish.ServiceRoot{
		OdataContext:   "/redfish/v1/$metadata#ServiceRoot.ServiceRoot",
		OdataID:        "/redfish/v1/",
		OdataType:      "#ServiceRoot.v1_17_0.ServiceRoot",
		ID:             "RootService",
		Name:           "Root Service",
		RedfishVersion: "1.17.0",
		UUID:           "123e4567-e89b-12d3-a456-426614174000",
		SessionService: struct {
			OdataID string `json:"@odata.id"`
		}{
			OdataID: "/redfish/v1/SessionService",
		},
		Systems: struct {
			OdataID string `json:"@odata.id"`
		}{
			OdataID: "/redfish/v1/Systems",
		},
		Links: struct {
			Sessions struct {
				OdataID string `json:"@odata.id"`
			} `json:"Sessions"`
		}{
			Sessions: struct {
				OdataID string `json:"@odata.id"`
			}{
				OdataID: "/redfish/v1/SessionService/Sessions",
			},
		},
	}
	c.JSON(http.StatusOK, serviceRoot)
}
func (r *deviceManagementRoutes) getSessionCollection(c *gin.Context) {
	sessionCollection := redfish.SessionCollection{
		OdataContext: "/redfish/v1/$metadata#SessionCollection.SessionCollection",
		OdataID:      "/redfish/v1/SessionService/Sessions",
		OdataType:    "#SessionCollection.SessionCollection",
		Name:         "Session Collection",
		Members: []struct {
			OdataID string `json:"@odata.id"`
		}{
			{OdataID: "/redfish/v1/SessionService/Sessions/1"},
		},
		MembersOdataCount: 1,
	}
	c.JSON(http.StatusOK, sessionCollection)
}
func (r *deviceManagementRoutes) getSessionService(c *gin.Context) {
	sessionService := redfish.SessionService{
		OdataContext: "/redfish/v1/$metadata#SessionService.SessionService",
		OdataID:      "/redfish/v1/SessionService",
		OdataType:    "#SessionService.SessionService",
		ID:           "SessionService",
		Name:         "Session Service",
	}
	c.JSON(http.StatusOK, sessionService)
}
func (r *deviceManagementRoutes) getSession(c *gin.Context) {
	session := redfish.Session{
		OdataContext: "/redfish/v1/$metadata#Session.Session",
		OdataID:      "/redfish/v1/SessionService/Sessions/1",
		OdataType:    "#Session.v1_0_2.Session",
		ID:           "1",
		Name:         "Session 1",
	}
	c.JSON(http.StatusOK, session)
}

func (r *deviceManagementRoutes) getComputerSystem(c *gin.Context) {
	computerSystem := redfish.ComputerSystem{
		OdataContext: "/redfish/v1/$metadata#ComputerSystem.ComputerSystem",
		OdataID:      "/redfish/v1/Systems/1",
		OdataType:    "#ComputerSystem.v1_22_0.ComputerSystem",
		ID:           "1",
		Name:         "Example Computer System",
		Manufacturer: "Example Manufacturer",
	}
	c.JSON(http.StatusOK, computerSystem)
}
func (r *deviceManagementRoutes) getComputerSystemCollection(c *gin.Context) {
	ComputerSystemCollection := redfish.ComputerSystemCollection{
		OdataContext: "/redfish/v1/$metadata#ComputerSystemCollection.ComputerSystemCollection",
		OdataID:      "/redfish/v1/Systems",
		OdataType:    "#ComputerSystemCollection.ComputerSystemCollection",
		Name:         "Systems Collection",
		Members: []struct {
			OdataID string `json:"@odata.id"`
		}{
			{OdataID: "/redfish/v1/Systems/1"},
		},
		MembersOdataCount: 1,
	}
	c.JSON(http.StatusOK, ComputerSystemCollection)
}
func (r *deviceManagementRoutes) getMetadata(c *gin.Context) {
	metadata2 := `<?xml version="1.0" encoding="UTF-8"?>
<edmx:Edmx Version="4.0" xmlns:edmx="http://docs.oasis-open.org/odata/ns/edmx">
    <edmx:Reference Uri="http://redfish.dmtf.org/schemas/v1/RedfishExtensions_v1.xml">
      <edmx:Include Namespace="RedfishExtensions.v1_0_0" Alias="Redfish"/>
    </edmx:Reference>
    <edmx:Reference Uri="http://redfish.dmtf.org/schemas/v1/ComputerSystemCollection_v1.xml">
        <edmx:Include Namespace="ComputerSystemCollection"/>
    </edmx:Reference>
    <edmx:Reference Uri="http://redfish.dmtf.org/schemas/v1/ComputerSystem_v1.xml">
        <edmx:Include Namespace="ComputerSystem"/>
        <edmx:Include Namespace="ComputerSystem.v1_22_1"/>
    </edmx:Reference>
</edmx:Edmx>`
	c.Data(http.StatusOK, "application/xml", []byte(metadata2))
}
