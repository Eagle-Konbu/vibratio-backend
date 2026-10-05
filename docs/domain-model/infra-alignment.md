# infra との整合

本設計に合わせて、`vibratio-infra` 側で必要になる変更（現状は 2026-10-05 時点の infra）。いずれも Domain Model には影響しない。テーブルのキー設計などの詳細は別途決める。

| 対象 | 現状 | 必要な変更 | 根拠 |
| --- | --- | --- | --- |
| Episode の保存先 | S3 の `episodes/` に保存する前提 | Episode 用の DynamoDB テーブルを新設する（config テーブルとは分ける）。バッチの書き込み権限と AppSync のデータソースを追加する | CMS から AppSync 経由で Episode を一覧できるようにするため。Episode は設定ではなく、再生成可能な生成データであり、保護の方針も権限も設定とは異なる |
| S3 の `episodes/` と `articles/` | バッチが読み書きできる | 削除する | Article は保存しない（[D-01](article.md#d-01)）。Episode は DynamoDB に保存する |
| README の Data Model と設計方針 | 「設定は DynamoDB、生成データは S3」 | 「音声は S3、それ以外は DynamoDB」に更新する | 上記の変更に合わせる |
| CMS の Source スキーマ | `type: RSS \| ATOM` を持つ | `type` を削除する | [D-10](source.md#d-10) |
| 音声の配信キャッシュ | CloudFront が音声を最大1日キャッシュする。再生成で同じ Key に上書きすると、古い音声が配信される | 上書き時にキャッシュを無効化するか、生成ごとに新しい Key にするかを決める（未決） | Key の作り方は AudioStorage Adapter の内部の問題であり、どちらを選んでも Domain は変わらない（[D-12](episode.md#d-12)） |
