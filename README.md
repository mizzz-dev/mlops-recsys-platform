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
- Runbook、ADR、architectureを含む運用前提のドキュメント

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

## 主要コマンド

```bash
make help
make test
make train
make run-api
make request-sample
make compile-pipeline
make docker-build
```
