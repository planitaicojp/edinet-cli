package cmdutil

import (
	"fmt"

	"github.com/planitaicojp/edinet-cli/internal/api"
	"github.com/planitaicojp/edinet-cli/internal/config"
	cerrors "github.com/planitaicojp/edinet-cli/internal/errors"
	"github.com/spf13/cobra"
)

func NewClient(cmd *cobra.Command) (*api.Client, error) {
	apiKey := config.ResolveAPIKey(GetStringFlag(cmd, "api-key"))
	if apiKey == "" {
		return nil, &cerrors.AuthError{
			Message: fmt.Sprintf("API キーが設定されていません。%s 環境変数を設定するか、'edinet config set api-key <key>' を実行してください", config.EnvAPIKey),
		}
	}

	client := api.NewClient(apiKey)
	client.Debug = GetBoolFlag(cmd, "verbose")
	return client, nil
}
