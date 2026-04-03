package config

import (
	"fmt"

	"github.com/planitaicojp/edinet-cli/internal/api"
	internalConfig "github.com/planitaicojp/edinet-cli/internal/config"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "現在の設定を表示",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := internalConfig.Load()
		if err != nil {
			return err
		}

		w := cmd.OutOrStdout()

		maskedKey := "(未設定)"
		if cfg.APIKey != "" {
			maskedKey = api.MaskAPIKey(cfg.APIKey)
		}

		_, _ = fmt.Fprintf(w, "config_dir:      %s\n", internalConfig.ConfigDir())
		_, _ = fmt.Fprintf(w, "api_key:         %s\n", maskedKey)
		_, _ = fmt.Fprintf(w, "default_format:  %s\n", cfg.DefaultFormat)

		resolved := internalConfig.ResolveAPIKey("")
		if resolved != "" && resolved != cfg.APIKey {
			_, _ = fmt.Fprintf(w, "\n(環境変数 %s により上書きされています)\n", internalConfig.EnvAPIKey)
		}

		return nil
	},
}
