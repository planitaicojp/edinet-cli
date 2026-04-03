package cmdutil

import (
	"os"

	"github.com/planitaicojp/edinet-cli/internal/config"
	"github.com/planitaicojp/edinet-cli/internal/output"
	"github.com/spf13/cobra"
)

func GetFormatter(cmd *cobra.Command) output.Formatter {
	format := GetStringFlag(cmd, "format")
	if format == "" {
		format = os.Getenv(config.EnvFormat)
	}
	if format == "" {
		cfg, err := config.Load()
		if err == nil && cfg.DefaultFormat != "" {
			format = cfg.DefaultFormat
		}
	}
	return output.New(format)
}
