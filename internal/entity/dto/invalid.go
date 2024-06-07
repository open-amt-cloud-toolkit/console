package dto

import "github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"

type NotValidError struct {
	Console consoleerrors.InternalError
}

func (e NotValidError) Error() string {
	return e.Console.Error()
}

func (e NotValidError) Wrap(call, function string, err error) error {
	_ = e.Console.Wrap(call, function, err)
	e.Console.Message = "Invalid input: " + err.Error()

	return e
}
