package config

import (
	"fmt"

	internalConfig "github.com/planitaicojp/edinet-cli/internal/config"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "設定値を変更",
	Example: `  edinet config set api-key YOUR_API_KEY
  edinet config set default-format json`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], args[1]

		cfg, err := internalConfig.Load()
		if err != nil {
			return err
		}

		switch key {
		case "api-key":
			cfg.APIKey = value
		case "default-format":
			cfg.DefaultFormat = value
		default:
			return fmt.Errorf("不明な設定キー: %s (api-key, default-format)", key)
		}

		if err := internalConfig.Save(cfg); err != nil {
			return err
		}
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "設定を保存しました: %s\n", key)
		return nil
	},
}
