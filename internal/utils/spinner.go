package utils

import (
	"time"

	"github.com/briandowns/spinner"
)

func NewSpinner(text string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + text
	s.Color("cyan")
	return s
}
