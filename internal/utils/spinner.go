package utils

import (
	"time"

	"github.com/briandowns/spinner"
)

func NewSpinner(msg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // ⠋ ⠙ ⠹ etc.
	s.Suffix = " " + msg
	return s
}
