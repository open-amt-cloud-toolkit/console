package profiles

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jritsema/go-htmx-starter/pkg/webtools"
	"github.com/jritsema/gotoolbox/web"
	"go.etcd.io/bbolt"
	"gopkg.in/yaml.v3"
)

type Profile struct {
	Id            int           `yaml:"id" json:"id"`
	Name          string        `yaml:"name" json:"name"`
	Technology    string        `yaml:"technology" json:"technology"`
	Configuration Configuration `yaml:"configuration" json:"configuration"`
}

type Configuration struct {
	RemoteManagement    RemoteManagement    `yaml:"remoteManagement" json:"remoteManagement"`
	EnterpriseAssistant EnterpriseAssistant `yaml:"enterpriseAssistant,omitempty" json:"enterpriseAssistant,omitempty"`
	AMTSpecific         AMTSpecific         `yaml:"amtSpecific,omitempty" json:"amtSpecific,omitempty"`
	BMCSpecific         BMCSpecific         `yaml:"bmcSpecific,omitempty" json:"bmcSpecific,omitempty"`
	DASHSpecific        DASHSpecific        `yaml:"dashSpecific,omitempty" json:"dashSpecific,omitempty"`
	RedfishSpecific     RedfishSpecific     `yaml:"redfishSpecific,omitempty" json:"redfishSpecific,omitempty"`
}

type RemoteManagement struct {
	GeneralSettings GeneralSettings        `yaml:"generalSettings,omitempty" json:"generalSettings,omitempty"`
	Network         Network                `yaml:"network,omitempty" json:"network,omitempty"`
	Authentication  AuthenticationProfiles `yaml:"authentication,omitempty" json:"authentication,omitempty"`
	TLS             TLS                    `yaml:"tls,omitempty" json:"tls,omitempty"`
	Redirection     Redirection            `yaml:"redirection,omitempty" json:"redirection,omitempty"`
	Accounts        []Account              `yaml:"accounts,omitempty" json:"accounts,omitempty"`
	AdminPassword   string                 `yaml:"adminPassword" json:"adminPassword" validate:"required,min=8,max=32"`
}

type GeneralSettings struct {
	HostName            string `yaml:"hostName,omitempty" json:"hostName,omitempty"`
	DomainName          string `yaml:"domainName,omitempty" json:"domainName,omitempty"`
	SharedFQDN          bool   `yaml:"sharedFQDN,omitempty" json:"sharedFQDN,omitempty"`
	NetworkEnabled      bool   `yaml:"networkEnabled,omitempty" json:"networkEnabled,omitempty"`
	PingResponseEnabled bool   `yaml:"pingResponseEnabled,omitempty" json:"pingResponseEnabled,omitempty"`
}

type Network struct {
	Wired    WiredNetwork           `yaml:"wired,omitempty" json:"wired,omitempty"`
	Wireless WirelessNetworkProfile `yaml:"wireless,omitempty" json:"wireless,omitempty"`
}

type WiredNetwork struct {
	DHCPEnabled               bool   `yaml:"dhcpEnabled,omitempty" json:"dhcpEnabled,omitempty"`
	IPSyncEnabled             bool   `yaml:"ipSyncEnabled,omitempty" json:"ipSyncEnabled,omitempty"`
	SharedStaticIP            bool   `yaml:"sharedStaticIP,omitempty" json:"sharedStaticIP,omitempty"`
	IPAddress                 string `yaml:"ipAddress,omitempty" json:"ipAddress,omitempty"`
	SubnetMask                string `yaml:"subnetMask,omitempty" json:"subnetMask,omitempty"`
	DefaultGateway            string `yaml:"defaultGateway,omitempty" json:"defaultGateway,omitempty"`
	PrimaryDNS                string `yaml:"primaryDNS,omitempty" json:"primaryDNS,omitempty"`
	SecondaryDNS              string `yaml:"secondaryDNS,omitempty" json:"secondaryDNS,omitempty"`
	AuthenticationProfileName string `yaml:"authenticationProfileName,omitempty" json:"authenticationProfileName,omitempty"`
}

type WirelessNetworkProfile struct {
	Profiles []WirelessNetwork `yaml:"profiles,omitempty" json:"profiles,omitempty" validate:"max=8"` // Limit of 8 profiles
}

type WirelessNetwork struct {
	ProfileName               string `yaml:"profileName,omitempty" json:"profileName,omitempty"`
	SSID                      string `yaml:"ssid,omitempty" json:"ssid,omitempty"`
	Priority                  int    `yaml:"priority,omitempty" json:"priority,omitempty"`
	WifiSecurity              string `yaml:"wifiSecurity,omitempty" json:"wifiSecurity,omitempty" validate:"oneof=WPA WPA2 WPA-Enterprise WPA2-Enterprise WPA3-SAE WPA3-OWE"`
	Passphrase                string `yaml:"passphrase,omitempty" json:"passphrase,omitempty"`
	AuthenticationProfileName string `yaml:"authenticationProfileName,omitempty" json:"authenticationProfileName,omitempty"`
}

type AuthenticationProfiles struct {
	Profiles []Authentication `yaml:"profiles,omitempty" json:"profiles,omitempty" validate:"max=7"` // Limit of 7 certificates in AMT certificate store
}

type Authentication struct {
	ProfileName            string `yaml:"profileName,omitempty" json:"profileName,omitempty"`
	Username               string `yaml:"username,omitempty" json:"username,omitempty"`
	Password               string `yaml:"password,omitempty" json:"password,omitempty"`
	AuthenticationProtocol string `yaml:"authenticationProtocol,omitempty" json:"authenticationProtocol,omitempty"`
	ClientCert             string `yaml:"clientCert,omitempty" json:"clientCert,omitempty"`
	CACert                 string `yaml:"caCert,omitempty" json:"caCert,omitempty"`
	PrivateKey             string `yaml:"privateKey,omitempty" json:"privateKey,omitempty"`
}

type TLS struct {
	MutualAuthentication bool     `yaml:"mutualAuthentication,omitempty" json:"mutualAuthentication,omitempty"`
	Enabled              bool     `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	TrustedCN            []string `yaml:"trustedCN,omitempty" json:"trustedCN,omitempty" validate:"max=10"` // Limit of 10 items
}

type Redirection struct {
	Enabled     bool                `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Services    RedirectionServices `yaml:"services,omitempty" json:"services,omitempty"`
	UserConsent string              `yaml:"userConsent,omitempty" json:"userConsent,omitempty" validate:"oneof=KVM All None"`
}

type RedirectionServices struct {
	KVM  bool `yaml:"kvm,omitempty" json:"kvm,omitempty"`
	SOL  bool `yaml:"sol,omitempty" json:"sol,omitempty"`
	IDER bool `yaml:"ider,omitempty" json:"ider,omitempty"`
}

type EnterpriseAssistant struct {
	URL      string `yaml:"url,omitempty" json:"url,omitempty"`
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
}

type AMTSpecific struct {
	ControlMode         string `yaml:"controlMode,omitempty" json:"controlMode,omitempty" validate:"oneof=acm ccm"`
	ProvisioningCert    string `yaml:"provisioningCert,omitempty" json:"provisioningCert,omitempty"`
	ProvisioningCertPwd string `yaml:"provisioningCertPwd,omitempty" json:"provisioningCertPwd,omitempty"`
	MEBXPassword        string `yaml:"mebxPassword,omitempty" json:"mebxPassword,omitempty" validate:"required,min=8,max=32"`
}

type Account struct {
	Credential Credential `yaml:"credential,omitempty" json:"credential,omitempty"`
	Scopes     []string   `yaml:"scopes,omitempty" json:"scopes,omitempty"`
}

type Credential struct {
	DigestUsername  string `yaml:"username,omitempty" json:"username,omitempty"`
	DigestPassword  string `yaml:"password,omitempty" json:"password,omitempty"`
	KerberosUserSid string `yaml:"kerberosSID,omitempty" json:"kerberosSID,omitempty"`
}

type BMCSpecific struct {
}

type DASHSpecific struct {
}

type RedfishSpecific struct {
}

func (p *Profile) IsValid() (bool, []string) {
	errors := []string{}
	if p.Name == "" {
		errors = append(errors, "Name is required")
	}
	if p.Configuration.RemoteManagement.AdminPassword == "" {
		errors = append(errors, "AMT Password is required")
	}
	return len(errors) <= 0, errors
}

// Get all profiles
func (pt ProfileThing) GetProfiles() []Profile {
	var data []Profile
	_ = pt.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Profiles"))
		_ = b.ForEach(func(k, v []byte) error {
			result := &Profile{}

			// Marshal user data into bytes.
			err := json.Unmarshal(v, result)
			if err != nil {
				return err
			}
			data = append(data, *result)
			return nil
		})
		return nil
	})
	return data
}

func (pt ProfileThing) GetProfileByID(id string) Profile {
	result := Profile{}

	_ = pt.db.Update(func(tx *bbolt.Tx) error {
		// Get buckets
		b := tx.Bucket([]byte("Profiles"))
		intId, _ := strconv.Atoi(id)
		profileSlice := b.Get(itob(intId))

		// Marshal user data into bytes.
		err := json.Unmarshal(profileSlice, &result)
		if err != nil {
			return err
		}
		return nil
	})
	return result
}

func (pt ProfileThing) UpdateProfile(profile Profile) {
	_ = pt.db.Update(func(tx *bbolt.Tx) error {
		// Get buckets
		b := tx.Bucket([]byte("Profiles"))

		profileSlice := b.Get(itob(profile.Id))
		result := &Profile{}

		// Marshal user data into bytes.
		err := json.Unmarshal(profileSlice, result)
		if err != nil {
			return err
		}
		result.Name = profile.Name

		// Marshal user data into bytes.
		buf, err := json.Marshal(profile)
		if err != nil {
			return err
		}
		// Persist bytes to users bucket.
		return b.Put(itob(profile.Id), buf)
	})
}

func (pt ProfileThing) AddProfile(profile Profile) error {
	err := pt.db.Update(func(tx *bbolt.Tx) error {
		// Get buckets
		b := tx.Bucket([]byte("Profiles"))

		// Generate ID for the profile.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := b.NextSequence()
		profile.Id = int(id)

		// Marshal user data into bytes.
		buf, err := json.Marshal(profile)
		if err != nil {
			return err
		}
		// Persist bytes to users bucket.
		return b.Put(itob(profile.Id), buf)
	})
	return err
}

func (pt ProfileThing) DeleteProfile(id string) {
	_ = pt.db.Update(func(tx *bbolt.Tx) error {
		// Get buckets
		b := tx.Bucket([]byte("Profiles"))
		intId, _ := strconv.Atoi(id)
		_ = b.Delete(itob(intId))
		return nil
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func (pt ProfileThing) ExportProfile(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	profile := pt.GetProfileByID(id)
	yamlData, err := yaml.Marshal(&profile)
	if err != nil {
		return webtools.HTML(r, http.StatusBadRequest, pt.html, "profiles/errors.html", nil, nil)
	}

	headers := web.Headers{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s.yaml"`, profile.Name),
	}
	return &web.Response{
		Status:      200,
		ContentType: "application/x-yaml",
		Content:     bytes.NewReader(yamlData),
		Headers:     headers,
	}
}

func (pt ProfileThing) Download(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	headers := web.Headers{
		"HX-Redirect": "http://localhost:8080/profile/export/" + id,
	}
	return &web.Response{
		Status:      200,
		ContentType: "",
		Content:     bytes.NewReader([]byte("test data")),
		Headers:     headers,
	}
}

// checkboxValue converts the value of a checkbox to a boolean
func checkboxValue(s string) bool {
	lowercase := strings.ToLower(s)
	return lowercase == "on"
}
