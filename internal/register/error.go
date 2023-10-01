package register

import (
	"errors"
	"fmt"

	"modernc.org/sqlite"
)

var (
	BackendStorageError = errors.New("")
	DataIntegrityError  = errors.New("")
)

func errorWrapper(message string, parentError error) error {
	if parentError == nil {
		return nil
	}

	if dbErr, ok := parentError.(*sqlite.Error); ok {
		return fmt.Errorf("register:%s - %s%w", message, dbErr.Error(), BackendStorageError)
	}

	return fmt.Errorf("register:%s - %s%w", message, parentError.Error(), BackendStorageError)
}
