package profiles

import (
	"encoding/binary"
	"encoding/json"
	"strconv"

	"go.etcd.io/bbolt"
)

type Profile struct {
	Id   int
	Name string
}

func (p *Profile) IsValid() (bool, []string) {
	errors := []string{}
	if p.Name == "" {
		errors = append(errors, "Name is required")
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
