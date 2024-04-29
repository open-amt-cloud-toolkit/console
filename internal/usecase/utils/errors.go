package utils

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrParseVersion = errors.New("failed to parse version")
)
