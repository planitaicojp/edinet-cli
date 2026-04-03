package document

import (
	"fmt"
	"os"
	"time"

	"github.com/planitaicojp/edinet-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "指定日付の書類一覧を取得",
	Long:  `指定した日付に提出された開示書類の一覧を取得します。`,
	Example: `  edinet document list --date 2024-04-01
  edinet document list --date today --type 2
  edinet document list --date yesterday --format json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		date, _ := cmd.Flags().GetString("date")
		date = resolveDate(date)

		typ, _ := cmd.Flags().GetInt("type")

		resp, err := client.ListDocuments(date, typ)
		if err != nil {
			return err
		}

		formatter := cmdutil.GetFormatter(cmd)

		if typ == 2 && len(resp.Results) > 0 {
			return formatter.Format(os.Stdout, resp.Results)
		}

		_, _ = fmt.Fprintf(os.Stdout, "日付: %s / 件数: %d / 更新日時: %s\n",
			resp.Metadata.Parameter.Date,
			resp.Metadata.ResultSet.Count,
			resp.Metadata.ProcessDateTime)
		return nil
	},
}

func init() {
	listCmd.Flags().String("date", "", "対象日付 (YYYY-MM-DD, today, yesterday) [必須]")
	listCmd.Flags().Int("type", 1, "取得情報 (1: メタデータのみ, 2: 提出書類一覧+メタデータ)")
	_ = listCmd.MarkFlagRequired("date")
}

func resolveDate(s string) string {
	now := time.Now()
	switch s {
	case "today":
		return now.Format("2006-01-02")
	case "yesterday":
		return now.AddDate(0, 0, -1).Format("2006-01-02")
	default:
		return s
	}
}
