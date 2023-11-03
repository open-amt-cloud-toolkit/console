package certificates

// import "strconv"

var data []Certificate

type Certificate struct {
	UUID      string
	Name      string
	IPAddress string
	FWVersion string
}

func init() {
	data = []Certificate{
		{
			UUID:      "12345678-abcd-efgh-ijkl-123456789001",
			Name:      "AMT Device 1",
			IPAddress: "192.168.0.1",
			FWVersion: "15.1.123",
		},
		{
			UUID:      "12345678-abcd-efgh-ijkl-123456789002",
			Name:      "AMT Device 2",
			IPAddress: "192.168.0.2",
			FWVersion: "16.0.43",
		},
		{
			UUID:      "12345678-abcd-efgh-ijkl-123456789003",
			Name:      "AMT Device 3",
			IPAddress: "192.168.0.3",
			FWVersion: "16.1.25",
		},
	}
}

func (dt CertificateThing) GetByID(id string) Certificate {
	var result Certificate
	for _, i := range data {
		if i.UUID == id {
			result = i
			break
		}
	}
	return result
}

func (dt CertificateThing) Update(device Certificate) {
	result := []Certificate{}
	for _, i := range data {
		if i.UUID == device.UUID {
			i.UUID = device.UUID
			i.Name = device.Name
			i.IPAddress = device.IPAddress
			i.FWVersion = device.FWVersion
		}
		result = append(result, i)
	}
	data = result
}

func (dt CertificateThing) Add(device Certificate) {
	// max := 0
	// for _, i := range data {
	// 	n, _ := strconv.Atoi(i.UUID)
	// 	if n > max {
	// 		max = n
	// 	}
	// }
	// max++
	// id := strconv.Itoa(max)

	data = append(data, Certificate{
		UUID:      device.UUID,
		Name:      device.Name,
		IPAddress: device.IPAddress,
		FWVersion: device.FWVersion,
	})
}

func (dt CertificateThing) Delete(id string) {
	result := []Certificate{}
	for _, i := range data {
		if i.UUID != id {
			result = append(result, i)
		}
	}
	data = result
}
