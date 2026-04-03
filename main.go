package main

import (
	"fmt"
	"os"

	"github.com/planitaicojp/edinet-cli/cmd"
	cerrors "github.com/planitaicojp/edinet-cli/internal/errors"
)

func main() {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(cerrors.GetExitCode(err))
	}
}
