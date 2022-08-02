package tools

import (
	"io"
	"fmt"
	"net"
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"github.com/AlecAivazis/survey/v2/terminal"
)

// FlagError is the kind of error raised in flag processing
type FlagError struct {
	Err error
}

func (fe FlagError) Error() string {
	return fe.Err.Error()
}

func (fe FlagError) Unwrap() error {
	return fe.Err
}

// SilentError is an error that triggers exit code 1 without any error messaging
var SilentError = errors.New("SilentError")

// CancelError signals user-initiated cancellation
var CancelError = errors.New("CancelError")

func IsUserCancellation(err error) bool {
	return errors.Is(err, CancelError) || errors.Is(err, terminal.InterruptErr)
}

func MutuallyExclusive(message string, conditions ...bool) error {
	numTrue := 0

	for _, ok := range conditions {
		if ok {
			numTrue++
		}
	}

	if numTrue > 1 {
		return &FlagError{Err: errors.New(message)}
	}

	return nil
}

func PrintError(out io.Writer, err error, cmd *cobra.Command, debug bool) {
	var dnsError *net.DNSError

	if errors.As(err, &dnsError) {
		fmt.Fprintf(out, "error connecting to %s\n", dnsError.Name)

		if debug {
			fmt.Fprintln(out, dnsError)
		}

		return
	}

	fmt.Fprintln(out, err)

	var flagError *FlagError
	if errors.As(err, &flagError) || strings.HasPrefix(err.Error(), "unknown command ") {
		if !strings.HasSuffix(err.Error(), "\n") {
			fmt.Fprintln(out)
		}

		fmt.Fprintln(out, cmd.UsageString())
	}
}
