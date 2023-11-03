package devices

import (
	"encoding/binary"
	"encoding/json"
	"strconv"

	"go.etcd.io/bbolt"
)

type Device struct {
	Id        int
	UUID      string
	Name      string
	Address   string
	FWVersion string
	Username  string
	Password  string
}

func init() {
	// data = []Device{
	// 	{
	// 		UUID:      1,
	// 		Name:      "AMT Device 1",
	// 		IPAddress: "192.168.0.1",
	// 		FWVersion: "15.1.123",
	// 	},
	// 	{
	// 		UUID:      2,
	// 		Name:      "AMT Device 2",
	// 		IPAddress: "192.168.0.2",
	// 		FWVersion: "16.0.43",
	// 	},
	// 	{
	// 		UUID:      3,
	// 		Name:      "AMT Device 3",
	// 		IPAddress: "192.168.0.3",
	// 		FWVersion: "16.1.25",
	// 	},
	// }
}

// Get all devices
func (dt DeviceThing) GetDevices() []Device {
	var data []Device
	dt.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Devices"))
		b.ForEach(func(k, v []byte) error {
			result := &Device{}

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

func (dt DeviceThing) GetDeviceByID(id string) Device {
	result := Device{}

	dt.db.Update(func(tx *bbolt.Tx) error {
		// Get buckets
		b := tx.Bucket([]byte("Devices"))
		intId, _ := strconv.Atoi(id)
		deviceSlice := b.Get(itob(intId))

		// Marshal user data into bytes.
		err := json.Unmarshal(deviceSlice, &result)
		if err != nil {
			return err
		}
		return nil
	})
	return result
}

func (dt DeviceThing) UpdateDevice(device Device) {
	dt.db.Update(func(tx *bbolt.Tx) error {
		// Get buckets
		b := tx.Bucket([]byte("Devices"))

		deviceSlice := b.Get(itob(device.Id))
		result := &Device{}

		// Marshal user data into bytes.
		err := json.Unmarshal(deviceSlice, result)
		if err != nil {
			return err
		}
		result.FWVersion = device.FWVersion
		result.Address = device.Address
		result.Name = device.Name

		// Marshal user data into bytes.
		buf, err := json.Marshal(device)
		if err != nil {
			return err
		}
		// Persist bytes to users bucket.
		return b.Put(itob(device.Id), buf)
	})
}

func (dt DeviceThing) AddDevice(device Device) {
	dt.db.Update(func(tx *bbolt.Tx) error {
		// Get buckets
		b := tx.Bucket([]byte("Devices"))

		// Generate ID for the device.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := b.NextSequence()
		device.Id = int(id)

		// Marshal user data into bytes.
		buf, err := json.Marshal(device)
		if err != nil {
			return err
		}
		// Persist bytes to users bucket.
		return b.Put(itob(device.Id), buf)
	})
}

func (dt DeviceThing) DeleteDevice(id string) {
	dt.db.Update(func(tx *bbolt.Tx) error {
		// Get buckets
		b := tx.Bucket([]byte("Devices"))
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
