package devices

// MinAMTVersion - minimum AMT version required for certain features in power capabilities.
var MinAMTVersion = 9

// UseCase -.
type UseCase struct {
	repo   Repository
	device Management
}

// New -.
func New(r Repository, d Management) *UseCase {
	return &UseCase{
		repo:   r,
		device: d,
	}
}
