# mlops-recsys-platform

推薦モデルを実サービスに組み込むことを想定した、MLOps学習・実践用のOSSポートフォリオです。

このリポジトリでは、合成ユーザー行動ログから推薦モデルを学習し、Go製の推論APIとして配信し、fallback・メトリクス・CI・運用ドキュメントまで含めて、機械学習システムを継続運用する前提の最小構成を実装します。

## 目的

機械学習モデルを「作って終わり」にせず、Web APIとして安全に提供し、継続的に評価・監視・改善できる状態にすることを目的にしています。

このMVPで示すことは以下です。

- Goによる低依存の推論API実装
- Pythonによる合成データ生成、特徴量生成、学習、評価、モデル保存
- モデル未ロード時にもAPIを落とさないfallback設計
- GitHub ActionsによるGo/Python/Docker/pipeline compileのCI
- k6による推薦APIの負荷試験
- Cloud Runを想定したデプロイ設計とTerraform雛形
- Runbook、障害訓練、ADR、architectureを含む運用前提のドキュメント

## MVPスコープ

### 必須

- `GET /healthz`
- `GET /readyz`
- `GET /metrics`
- `GET /v1/recommendations?user_id={user_id}&limit={limit}`
- 合成イベントログ生成
- 簡易推薦モデルの学習・評価・保存
- モデルartifactの読み込み
- モデル未ロード時のfallback推薦
- ローカル実行用Makefile
- CI
- smoke負荷試験

### 対象外

- 本格的なFeature Store
- GKE本番運用
- A/Bテスト基盤
- 高度な深層学習モデル
- 管理画面
- 実ユーザーデータの利用

## ディレクトリ構成

```text
apps/inference-api/   Go製の推論API
ml/                   Python製のデータ生成・学習・評価処理
pipelines/training/   学習パイプライン定義とcompileスクリプト
data/                 サンプルデータとスキーマ
loadtests/k6/         負荷試験スクリプト
infra/terraform/      Cloud Run想定のTerraform雛形
docs/                 要件、設計、運用、ADR
.github/workflows/    CI
```

## ローカル実行

```bash
make test
make train
make run-api
make request-sample
```

`make train` で `artifacts/model.json` が生成されます。

## API例

```bash
curl 'http://localhost:8080/v1/recommendations?user_id=user_001&limit=3'
```

モデルが存在しない場合もAPIは落ちず、`strategy: fallback_popular` で返します。

## 負荷試験

APIを起動した状態で以下を実行します。

```bash
make loadtest
```

別URLに対して実行する場合は `BASE_URL` を指定します。

```bash
make loadtest BASE_URL=https://example.run.app
```

smoke負荷試験では以下を確認します。

- HTTP error rate < 1%
- p95 latency < 300ms
- 推薦レスポンスに `strategy` と `recommendations` が含まれること

## Cloud Runデプロイ設計

MVPでは実デプロイまでは対象外とし、`infra/terraform` にCloud Runを想定した最小構成の雛形を置いています。

詳細は以下を参照してください。

- `docs/deployment.md`
- `docs/gcp-setup.md`
- `docs/staging-release.md`
- `docs/staging-smoke-test.md`
- `docs/artifact-registry.md`
- `docs/rollback.md`

## 主要コマンド

```bash
make help
make test
make train
make run-api
make request-sample
make compile-pipeline
make docker-build
make loadtest
```
