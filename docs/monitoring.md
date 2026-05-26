# 監視設計

## APIメトリクス

- `recsys_requests_total`
- `recsys_errors_total`
- `recsys_fallback_total`
- `recsys_model_loaded`

## 監視したい指標

- p95 latency
- error rate
- fallback rate
- model version
- prediction distribution
- training data freshness

## アラート例

- fallback rateが通常より高い
- error rateが1%を超える
- model_loadedが0の状態が継続する
- p95 latencyが300msを超える
