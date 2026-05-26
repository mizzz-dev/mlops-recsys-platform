# デプロイ設計

## 目的

推論APIをCloud Runへデプロイする前提の構成を整理する。MVPでは実デプロイまでは行わず、Terraform雛形と運用上の確認観点を追加する。

## 前提

- 推論APIはコンテナとしてビルドする
- model artifactは将来的にGCSから取得する
- MVPではコンテナ起動時に `MODEL_PATH` からモデルを読む
- Secretや認証設定はproduction投入前に別途設計する

## 想定構成

```text
GitHub Actions
  -> Docker build
  -> Artifact Registry push
  -> Cloud Run deploy
  -> smoke test

Cloud Run
  -> inference-api container
  -> /healthz
  -> /readyz
  -> /metrics
  -> /v1/recommendations
```

## 環境変数

| 名前 | 内容 | MVPでの扱い |
|---|---|---|
| `PORT` | API listen port | `8080` |
| `MODEL_PATH` | モデルartifactのパス | `/app/artifacts/model.json` |

## デプロイ前チェック

- `make test` が通る
- `make docker-build` が通る
- `/healthz` が200を返す
- `/readyz` がモデルロード状態を返す
- モデル未ロード時にfallbackが返る
- k6 smoke負荷試験でp95 latencyが300ms未満

## production前の未対応事項

- Artifact Registry push
- Cloud Run deploy workflow
- GCSからのmodel artifact取得
- Cloud Run IAMまたはAPI認証
- Secret管理
- revision rollback手順の自動化
