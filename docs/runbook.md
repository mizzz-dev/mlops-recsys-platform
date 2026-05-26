# Runbook

## 推薦APIのエラー率が上がった場合

1. `/healthz` でプロセス状態を確認する
2. `/readyz` でモデルロード状態を確認する
3. `/metrics` で `recsys_errors_total` と `recsys_fallback_total` を確認する
4. 直近のデプロイ差分を確認する
5. モデルartifactの有無とJSON形式を確認する

## fallback rateが上がった場合

1. `recsys_model_loaded` が0か確認する
2. `MODEL_PATH` が正しいか確認する
3. `artifacts/model.json` が存在するか確認する
4. モデル生成ジョブの失敗有無を確認する
5. 必要に応じて旧モデルへ戻す

## モデル学習が失敗した場合

1. 入力データのスキーマを確認する
2. 合成データ生成が成功しているか確認する
3. 評価指標計算で例外が起きていないか確認する
4. artifact保存先の権限を確認する

## ロールバック方針

MVPではモデルartifactを再生成する。将来的にはmodel versionごとにartifactを保存し、API起動時または設定変更で旧モデルへ戻せるようにする。
