package errlst

import "github.com/pkg/errors"

var (
	AlreadyExists   = errors.New("User has already started bot")
	ConvertionError = errors.New("Impossible to convert interface")
	NilUpdateFields = errors.New("Nothing to update or Updating query undefined")
	ServerError     = errors.New("Server Error")
)
