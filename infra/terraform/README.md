# Terraform skeleton

## 目的

Cloud Runへ推論APIをデプロイするための最小Terraform雛形です。

MVPでは実デプロイまでは行わず、以下をレビューできる状態に留めます。

- Cloud Run service
- Artifact Registry imageを受け取る変数
- 環境変数 `MODEL_PATH`
- 環境変数 `MODEL_URI`
- GCS model artifact bucket
- Cloud Run runtime service accountへの最小IAM方針
- 最小限のリソース制限

## 初期化

```bash
cd infra/terraform
terraform init
terraform fmt -check
terraform validate
```

## Model artifact bucket

このTerraform雛形では、モデルartifact保存用のGCS bucketを作成します。

- bucket versioningを有効化
- uniform bucket-level accessを有効化
- 古いobject versionを一定数で削除するlifecycle ruleを設定
- `cloud_run_runtime_service_account` が指定された場合のみ `roles/storage.objectViewer` を付与

## 注意

- `image` にはArtifact Registry上のコンテナイメージを指定します。
- `MODEL_URI` はTerraform変数として渡せますが、API側のGCS download実装は後続対応です。
- production公開前にCloud Run IAM、認証、Secret管理、state backendを設計してください。
- state backendは未設定です。個人検証ではlocal state、チーム運用ではGCS backendを利用してください。
