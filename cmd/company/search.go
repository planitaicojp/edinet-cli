package company

import (
	"fmt"
	"os"

	"github.com/planitaicojp/edinet-cli/cmd/cmdutil"
	"github.com/planitaicojp/edinet-cli/internal/api"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "企業を検索",
	Long:  `ローカルキャッシュのEDINETコードリストから企業を部分一致検索します。`,
	Example: `  edinet company search トヨタ
  edinet company search 7203 --by code
  edinet company search E00001 --by edinet-code
  edinet company search toyota --by name`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		by, _ := cmd.Flags().GetString("by")
		cachePath := api.CacheFilePath()

		companies, err := api.LoadCodeList(cachePath)
		if err != nil {
			return err
		}

		results := api.SearchCompanies(companies, query, by)
		if len(results) == 0 {
			fmt.Fprintf(os.Stderr, "検索結果: 0件 (query=%q, by=%s)\n", query, by)
			return nil
		}

		fmt.Fprintf(os.Stderr, "検索結果: %d件\n", len(results))
		formatter := cmdutil.GetFormatter(cmd)
		return formatter.Format(os.Stdout, results)
	},
}

func init() {
	searchCmd.Flags().String("by", "all", "検索対象 (name, code, edinet-code, all)")
}
