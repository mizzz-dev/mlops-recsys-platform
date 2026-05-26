# 障害訓練

## 目的

推薦APIの運用時に起こりやすい障害を想定し、検知・切り分け・復旧までの流れを確認できるようにする。

## 訓練1: モデルartifactが存在しない

### 想定

`artifacts/model.json` が存在しない、または `MODEL_PATH` が誤っている。

### 期待される挙動

- APIプロセスは起動する
- `/readyz` は `model_loaded: false` を返す
- `/v1/recommendations` は `strategy: fallback_popular` を返す
- `/metrics` の `recsys_model_loaded` は `0` になる

### 確認コマンド

```bash
rm -f artifacts/model.json
make run-api
curl -s 'http://localhost:8080/readyz' | python3 -m json.tool
curl -s 'http://localhost:8080/v1/recommendations?user_id=user_001&limit=3' | python3 -m json.tool
curl -s 'http://localhost:8080/metrics'
```

### 復旧

```bash
make train
make run-api
```

## 訓練2: fallback rateが上昇する

### 想定

モデルが読み込めない状態でリクエストが継続している。

### 確認観点

- `recsys_fallback_total` が増加しているか
- `recsys_model_loaded` が0か
- APIログにモデル読み込み失敗のwarningが出ているか

### 復旧

- `MODEL_PATH` を確認する
- artifactを再生成する
- Cloud Run想定では前revisionへ戻す

## 訓練3: API error rateが上昇する

### 想定

不正な `limit` や `user_id` 未指定のリクエストが増えている。

### 期待される挙動

- 不正入力は400で返る
- `recsys_errors_total` が増加する
- 500ではなくクライアントエラーとして扱える

### 確認コマンド

```bash
curl -i 'http://localhost:8080/v1/recommendations?limit=100'
curl -s 'http://localhost:8080/metrics'
```

## 訓練4: smoke負荷試験

### 想定

APIが最低限の負荷に耐えられるかを確認する。

### 確認コマンド

```bash
make train
make run-api
make loadtest
```

### 合格条件

- HTTP error rate < 1%
- p95 latency < 300ms
- レスポンスに `strategy` と `recommendations` が含まれる
