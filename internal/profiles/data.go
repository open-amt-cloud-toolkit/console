package profiles

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jritsema/go-htmx-starter/pkg/webtools"
	"github.com/jritsema/gotoolbox/web"
	"go.etcd.io/bbolt"
	"gopkg.in/yaml.v3"
)

type Profile struct {
	Id              int                `yaml:"id" json:"id"`
	Name            string             `yaml:"name" json:"name"`
	ControlMode     string             `yaml:"controlMode" json:"controlMode"`
	MEBXPassword    string             `yaml:"mebxPassword" json:"mebxPassword"`
	Activate        Activate           `yaml:"activate" json:"activate"`
	Network         Network            `yaml:"network" json:"network"`
	WirelessNetwork [7]WirelessNetwork `yaml:"wifiConfigs" json:"wifiConfigs"`
	IEEE8021xConfig [7]IEEE8021xConfig `yaml:"ieee8021xConfigs" json:"ieee8021xConfigs"`
}

type Activate struct {
	AMTPassword         string `yaml:"amtPassword" json:"amtPassword"`
	ProvisioningCert    string `yaml:"provisioningCert" json:"provisioningCert"`
	ProvisioningCertPwd string `yaml:"provisioningCertPwd" json:"provisioningCertPwd"`
}

type Network struct {
	DHCPSync bool
}

type WirelessNetwork struct {
	ProfileName string `yaml:"profileName"`
}

type IEEE8021xConfig struct {
}

func (p *Profile) IsValid() (bool, []string) {
	errors := []string{}
	if p.Name == "" {
		errors = append(errors, "Name is required")
	}
	if p.Activate.AMTPassword == "" {
		errors = append(errors, "AMT Password is required")
	}
	return len(errors) <= 0, errors
}

// Get all profiles
func (pt ProfileThing) GetProfiles() []Profile {
	var data []Profile
	pt.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Profiles"))
		b.ForEach(func(k, v []byte) error {
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

	pt.db.Update(func(tx *bbolt.Tx) error {
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
	pt.db.Update(func(tx *bbolt.Tx) error {
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

func (pt ProfileThing) AddDevice(profile Profile) {
	pt.db.Update(func(tx *bbolt.Tx) error {
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
}

func (pt ProfileThing) DeleteProfile(id string) {
	pt.db.Update(func(tx *bbolt.Tx) error {
		// Get buckets
		b := tx.Bucket([]byte("Profiles"))
		intId, _ := strconv.Atoi(id)
		b.Delete(itob(intId))
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
		"Content-Disposition": `attachment; filename="profile.yaml"`,
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
