package util

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// Die prints an error with its cause if given and exits.
func Die(msg string, cause error) {
	if cause != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(cause, msg).Error())
	} else {
		fmt.Fprintln(os.Stderr, msg)
	}
	os.Exit(1)
}

// PrintErrors is used to print several errors at once.
func PrintErrors(errs []error) {
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
