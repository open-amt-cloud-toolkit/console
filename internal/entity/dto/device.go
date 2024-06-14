package dto

import (
	"time"
)

type Device struct {
	ConnectionStatus bool        `json:"connectionStatus"`
	MPSInstance      string      `json:"mpsInstance"`
	Hostname         string      `json:"hostname"`
	GUID             string      `json:"guid"`
	MPSUsername      string      `json:"mpsusername"`
	Tags             []string    `json:"tags"`
	TenantID         string      `json:"tenantId"`
	FriendlyName     string      `json:"friendlyName"`
	DNSSuffix        string      `json:"dnsSuffix"`
	LastConnected    *time.Time  `json:"lastConnected,omitempty"`
	LastSeen         *time.Time  `json:"lastSeen,omitempty"`
	LastDisconnected *time.Time  `json:"lastDisconnected,omitempty"`
	DeviceInfo       *DeviceInfo `json:"deviceInfo,omitempty"`
	Username         string      `json:"username"`
	Password         string      `json:"password"`
	UseTLS           bool        `json:"useTLS"`
	AllowSelfSigned  bool        `json:"allowSelfSigned"`
}

type DeviceInfo struct {
	FWVersion   string    `json:"fwVersion"`
	FWBuild     string    `json:"fwBuild"`
	FWSku       string    `json:"fwSku"`
	CurrentMode string    `json:"currentMode"`
	Features    string    `json:"features"`
	IPAddress   string    `json:"ipAddress"`
	LastUpdated time.Time `json:"lastUpdated"`
}
