# 要件定義

## 目的

推薦モデルを実サービスに組み込むことを想定し、推論API・学習処理・CI・運用ドキュメントを備えたOSSポートフォリオを構築する。

## 想定ユースケース

- ユーザー行動ログをもとにおすすめコンテンツを返す
- モデル未ロードやモデル破損時もfallback推薦でサービス提供を継続する
- CIで推論API・学習処理・パイプライン定義を検証する
- READMEとRunbookから第三者が設計意図と運用手順を理解できる

## MVP必須要件

- Go製推論APIを提供する
- Pythonで合成データ生成、特徴量生成、学習、評価、モデル保存を行う
- 推薦APIはモデル利用時とfallback時の両方に対応する
- `/metrics` で最低限の運用メトリクスを公開する
- `make test` とCIで主要テストが通る
- 実データや個人情報を含めない

## 対象外

- 本格Feature Store
- GKE本番運用
- A/Bテスト基盤
- 高度な深層学習モデル
- 管理画面

## 受け入れ条件

- `make train` で `artifacts/model.json` が生成される
- `make run-api` でAPIが起動する
- 推薦APIが `strategy: model` または `strategy: fallback_popular` を返す
- モデルが存在しなくてもAPIが500で落ちない
- Go/Pythonのテストが通る
