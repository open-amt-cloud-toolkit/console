package entity

import "time"

type Device struct {
	ConnectionStatus bool
	MpsInstance      string
	Hostname         string
	GUID             string
	Mpsusername      string
	Tags             string
	TenantID         string
	FriendlyName     string
	DNSSuffix        string
	LastConnected    *time.Time
	LastSeen         *time.Time
	LastDisconnected *time.Time
	DeviceInfo       string
	Username         string
	Password         string
	UseTLS           bool
	AllowSelfSigned  bool
}
