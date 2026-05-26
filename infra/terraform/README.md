# Terraform skeleton

## 目的

Cloud Runへ推論APIをデプロイするための最小Terraform雛形です。

MVPでは実デプロイまでは行わず、以下をレビューできる状態に留めます。

- Cloud Run service
- Artifact Registry imageを受け取る変数
- 環境変数 `MODEL_PATH`
- 最小限のリソース制限
- 将来のIAM/Secret/GCS連携のTODO

## 初期化

```bash
cd infra/terraform
t terraform init
terraform fmt -check
terraform validate
```

## 注意

- `image` にはArtifact Registry上のコンテナイメージを指定します。
- production公開前にCloud Run IAM、認証、GCS artifact取得、Secret管理を追加してください。
- state backendは未設定です。個人検証ではlocal state、チーム運用ではGCS backendを利用してください。
