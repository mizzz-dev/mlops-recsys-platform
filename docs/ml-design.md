# ML設計

## 問題設定

ユーザー行動ログから人気コンテンツを推定し、推薦APIで返す。MVPでは高度な個別最適化より、学習から推論API提供までのライフサイクルを示すことを優先する。

## 入力データ

合成イベントログのみを扱う。

| カラム | 内容 |
|---|---|
| user_id | 合成ユーザーID |
| content_id | 合成コンテンツID |
| event_type | view / start / complete / like |
| timestamp | イベント時刻 |

## モデル

MVPでは `popular_baseline` を採用する。イベント種別に応じて重み付けし、人気順に推薦候補を返す。

## 評価指標

- Precision@10
- event_count
- user_count

## 採用理由

深いモデル精度よりも、API化、fallback、CI、運用設計を優先するため。高度なモデルは、基盤が安定した後の拡張対象とする。
