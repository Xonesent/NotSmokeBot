package errlst

import "github.com/pkg/errors"

var (
	NotStartedError      = errors.New("User did not start bot")
	AlreadyExistsError   = errors.New("User has already started bot")
	ConversionError      = errors.New("Impossible to convert interface")
	NilUpdateFieldsError = errors.New("Nothing to update or Updating query undefined")
	NothingFoundError    = errors.New("During find request nothing found")
	ServerError          = errors.New("Server Error")
	BlockError           = errors.New("forbidden, Forbidden: bot was blocked by the user")
)
