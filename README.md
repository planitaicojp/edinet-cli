# edinet-cli

金融庁が提供する [EDINET API v2](https://disclosure2dl.edinet-fsa.go.jp/guide/static/disclosure/WZEK0110.html) を操作するための CLI ツールです。
有価証券報告書等の開示書類を検索・取得できます。

## インストール

### Homebrew (macOS / Linux)

```bash
brew install planitaicojp/tap/edinet
```

### Scoop (Windows)

```powershell
scoop bucket add planitaicojp https://github.com/planitaicojp/bucket
scoop install edinet
```

### ソースからビルド

```bash
go install github.com/planitaicojp/edinet-cli@latest
```

### リリースバイナリ

[Releases](https://github.com/planitaicojp/edinet-cli/releases) ページからダウンロード：

**Linux (amd64)**

```bash
VERSION=$(curl -s https://api.github.com/repos/planitaicojp/edinet-cli/releases/latest | grep tag_name | cut -d '"' -f4)
curl -Lo edinet.tar.gz "https://github.com/planitaicojp/edinet-cli/releases/download/${VERSION}/edinet-cli_${VERSION#v}_linux_amd64.tar.gz"
tar xzf edinet.tar.gz edinet
sudo mv edinet /usr/local/bin/
rm edinet.tar.gz
```

**macOS (Apple Silicon)**

```bash
VERSION=$(curl -s https://api.github.com/repos/planitaicojp/edinet-cli/releases/latest | grep tag_name | cut -d '"' -f4)
curl -Lo edinet.tar.gz "https://github.com/planitaicojp/edinet-cli/releases/download/${VERSION}/edinet-cli_${VERSION#v}_darwin_arm64.tar.gz"
tar xzf edinet.tar.gz edinet
sudo mv edinet /usr/local/bin/
rm edinet.tar.gz
```

**Windows (amd64)**

```powershell
$version = (Invoke-RestMethod https://api.github.com/repos/planitaicojp/edinet-cli/releases/latest).tag_name
$v = $version -replace '^v', ''
Invoke-WebRequest -Uri "https://github.com/planitaicojp/edinet-cli/releases/download/$version/edinet-cli_${v}_windows_amd64.zip" -OutFile edinet.zip
Expand-Archive edinet.zip -DestinationPath .
Remove-Item edinet.zip
```

## セットアップ

### API キーの取得

1. [EDINET API](https://api.edinet-fsa.go.jp/api/auth/index.aspx?mode=1) でアカウントを作成
2. API キーを発行

### API キーの設定

```bash
# 方法1: 環境変数 (推奨)
export EDINET_API_KEY="your-api-key"

# 方法2: config ファイル
edinet config set api-key "your-api-key"
```

## 使い方

### 書類一覧の取得

```bash
# メタデータの取得
edinet document list --date 2024-04-01

# 提出書類一覧の取得
edinet document list --date today --type 2

# JSON形式で出力
edinet document list --date yesterday --type 2 --format json
```

### 書類のダウンロード

```bash
# XBRL (デフォルト)
edinet document download S100ABCD

# PDF
edinet document download S100ABCD --type pdf

# CSV
edinet document download S100ABCD --type csv --output ./data/
```

### 設定の確認

```bash
edinet config show
```

## 開発

```bash
# ビルド
make build

# テスト
make test

# リント
make lint

# 全部
make all
```

## 環境変数

| 変数名 | 説明 | デフォルト |
|--------|------|-----------|
| `EDINET_API_KEY` | API キー | - |
| `EDINET_FORMAT` | 出力形式 (table, json, csv) | `table` |
| `EDINET_CONFIG_DIR` | 設定ディレクトリ | `~/.config/edinet` |
| `EDINET_DEBUG` | デバッグログ | `false` |

## ライセンス

[Apache License 2.0](LICENSE)
