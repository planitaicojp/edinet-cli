package document

import (
	"fmt"
	"os"

	"github.com/planitaicojp/edinet-cli/cmd/cmdutil"
	cerrors "github.com/planitaicojp/edinet-cli/internal/errors"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <docID>",
	Short: "書類の詳細情報を表示",
	Long:  `指定した書類管理番号の詳細情報を表示します。書類一覧API (type=2) から該当書類を検索します。`,
	Example: `  edinet document show S100ABCD --date 2024-04-01
  edinet document show S100ABCD --date today`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		docID := args[0]
		date, _ := cmd.Flags().GetString("date")
		date = resolveDate(date)

		resp, err := client.ListDocuments(date, 2)
		if err != nil {
			return err
		}

		for _, doc := range resp.Results {
			if doc.DocID == docID {
				formatter := cmdutil.GetFormatter(cmd)
				return formatter.Format(os.Stdout, doc)
			}
		}

		return &cerrors.NotFoundError{
			Resource: "document",
			ID:       fmt.Sprintf("%s (date: %s)", docID, date),
		}
	},
}

func init() {
	showCmd.Flags().String("date", "", "対象日付 (YYYY-MM-DD, today, yesterday) [必須]")
	_ = showCmd.MarkFlagRequired("date")
}
