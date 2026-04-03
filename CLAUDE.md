# edinet-cli

## 概要

EDINET（Electronic Disclosure for Investors' NETwork）APIを操作するためのCLIツール。金融庁が提供する有価証券報告書等の開示書類を検索・取得する。

## 対象API

- **提供元**: 金融庁
- **API**: EDINET API v2 (https://disclosure2dl.edinet-fsa.go.jp/guide/static/disclosure/WZEK0110.html)
- **認証**: APIキーが必要（無料登録）
- **形式**: REST/JSON

## 主な機能

- 書類一覧の取得（日付指定、書類種別フィルタ）
- 書類の検索（企業名、EDINETコード、証券コードなど）
- 書類のダウンロード（XBRL、PDF、CSV、英文開示）
- 企業情報（EDINETコード一覧）の取得
- 有価証券報告書・四半期報告書・大量保有報告書等の種別対応

## 既存ツールの状況

- GitHub上にCLIが2つ存在するが、いずれも0スター・メンテナンス停止状態
- 成熟したCLI/SDKは存在しない

## 開発方針

- Go言語で実装
- サブコマンド構成（`edinet list`, `edinet search`, `edinet download` など）
- 出力形式: JSON（デフォルト）、テーブル表示
- XBRL解析の基本サポート
- Claude Code等のAIエージェントからの利用を想定した設計
- 日本語ドキュメント・ヘルプメッセージ
