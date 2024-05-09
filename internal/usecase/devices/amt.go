package devices

import "github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"

type AMTError struct {
	Console consoleerrors.InternalError
}

func (e AMTError) Error() string {
	return "amt error"
}
