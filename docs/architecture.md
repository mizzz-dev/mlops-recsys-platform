# アーキテクチャ

## 全体方針

MVPでは、ローカルで完走できる単純な構成を優先する。推論API、学習処理、パイプライン定義、運用ドキュメントを分離し、後からCloud Run、BigQuery、GCS、Vertex AI Pipelinesへ拡張しやすい構造にする。

## 構成

```text
[Python training]
  -> synthetic events
  -> feature build
  -> train popular baseline
  -> evaluate
  -> artifacts/model.json

[Go inference API]
  -> load artifacts/model.json
  -> /v1/recommendations
  -> /metrics
  -> fallback when model is missing

[GitHub Actions]
  -> Go test
  -> Python test
  -> pipeline compile
  -> Docker build
```

## 推論API

Go APIは起動時に `MODEL_PATH` からモデルartifactを読み込む。読み込みに失敗してもプロセスは終了せず、fallback推薦に切り替える。

## fallback設計

モデル未ロード、未知ユーザー、特徴量取得失敗時にAPIを落とすとユーザー影響が大きい。そのためMVPでは人気コンテンツの固定リストをfallbackとして返す。

## モデル更新

MVPでは起動時ロードのみ。将来的には以下を検討する。

- artifactのatomic切り替え
- model_version指定
- canary model routing
- rollback用の旧モデル保持
