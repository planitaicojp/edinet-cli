package document

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/planitaicojp/edinet-cli/cmd/cmdutil"
	cerrors "github.com/planitaicojp/edinet-cli/internal/errors"
	"github.com/spf13/cobra"
)

var docTypeMap = map[string]int{
	"xbrl":    1,
	"pdf":     2,
	"attach":  3,
	"english": 4,
	"csv":     5,
}

var docTypeExt = map[string]string{
	"xbrl":    ".zip",
	"pdf":     ".pdf",
	"attach":  ".zip",
	"english": ".zip",
	"csv":     ".zip",
}

var downloadCmd = &cobra.Command{
	Use:   "download <docID>",
	Short: "書類をダウンロード",
	Long:  `指定した書類管理番号の書類をダウンロードします。`,
	Example: `  edinet document download S100ABCD
  edinet document download S100ABCD --type pdf
  edinet document download S100ABCD --type csv --output ./data/`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		docID := args[0]
		typeName, _ := cmd.Flags().GetString("type")
		outputDir, _ := cmd.Flags().GetString("output")

		apiType, ok := docTypeMap[typeName]
		if !ok {
			return &cerrors.ValidationError{
				Field:   "type",
				Message: fmt.Sprintf("不正な書類タイプ: %s (xbrl, pdf, attach, english, csv)", typeName),
			}
		}

		body, _, err := client.DownloadDocument(docID, apiType)
		if err != nil {
			return err
		}
		defer func() { _ = body.Close() }()

		ext := docTypeExt[typeName]
		filename := fmt.Sprintf("%s_%s%s", docID, typeName, ext)
		outPath := filepath.Join(outputDir, filename)

		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("ファイル作成エラー: %w", err)
		}

		n, err := io.Copy(f, body)
		if err != nil {
			return fmt.Errorf("ダウンロードエラー: %w", err)
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("ファイル書き込みエラー: %w", err)
		}

		fmt.Fprintf(os.Stderr, "ダウンロード完了: %s (%d bytes)\n", outPath, n)
		return nil
	},
}

func init() {
	downloadCmd.Flags().StringP("type", "t", "xbrl", "書類タイプ (xbrl, pdf, attach, english, csv)")
	downloadCmd.Flags().StringP("output", "o", ".", "保存先ディレクトリ")
}
