package revip

import (
	"fmt"
)

// ErrFileNotFound should be returned if configuration file was not found.
type ErrFileNotFound struct {
	Path string
	Orig error
}

func (e *ErrFileNotFound) Error() string {
	return fmt.Sprintf("no such file: %q", e.Path)
}

//

// ErrPathNotFound should be returned if key (path) was not found in configuration.
type ErrPathNotFound struct {
	Path string
}

func (e *ErrPathNotFound) Error() string {
	return fmt.Sprintf("no key matched for path: %q", e.Path)
}
