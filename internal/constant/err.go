package constant

import "errors"

var (
	ErrInvalidParam  = errors.New(`invalid parameter`)
	ErrAuthorization = errors.New(`authorization error`)
)

var (
	ErrInternal = errors.New(`internal error`)
	ErrNotFound = errors.New(`not found`)
)

var (
	ErrInvalidFilePath = errors.New(`file path error`)
	ErrCompile         = errors.New(`compile error`)
	ErrCreateFile      = errors.New(`create file error`)
)
