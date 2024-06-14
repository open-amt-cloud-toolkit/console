package sqldb

import "github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"

type NotFoundError struct {
	Console consoleerrors.InternalError
}

func (e NotFoundError) Error() string {
	return e.Console.Error()
}

func (e NotFoundError) Wrap(call, function string, err error) error {
	_ = e.Console.Wrap(call, function, err)
	e.Console.Message = "requested resource not found"

	return e
}
