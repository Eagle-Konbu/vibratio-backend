# infra との整合

本設計に合わせて、`vibratio-infra` 側で必要になる変更（現状は 2026-10-05 時点の infra）。いずれも Domain Model には影響しない。Episode テーブルのキーは [architecture.md A-22](../architecture.md#a-22) のとおりとする。

| 対象 | 現状 | 必要な変更 | 根拠 |
| --- | --- | --- | --- |
| Episode の保存先 | S3 の `episodes/` に保存する前提 | Episode 用の DynamoDB テーブルを新設する（config テーブルとは分ける）。パーティションキーは Date とし、ソートキーは持たない。バッチの書き込み権限と AppSync のデータソースを追加する | CMS から AppSync 経由で Episode を一覧できるようにするため。Episode は設定ではなく、再生成可能な生成データであり、保護の方針も権限も設定とは異なる |
| S3 の `episodes/` と `articles/` | バッチが読み書きできる | 削除する | Article は保存しない（[D-01](article.md#d-01)）。Episode は DynamoDB に保存する |
| README の Data Model と設計方針 | 「設定は DynamoDB、生成データは S3」 | 「音声は S3、それ以外は DynamoDB」に更新する | 上記の変更に合わせる |
| CMS の Source スキーマ | `type: RSS \| ATOM` を持つ | `type` を削除する | [D-10](source.md#d-10) |
| generate-episode の Lambda のタイムアウト | 未確認 | 上限の 900 秒にする | 数十分の音声を1回の実行で生成する（[architecture.md A-21](../architecture.md#a-21)） |
| 音声の配信キャッシュ | CloudFront が音声を最大1日キャッシュする | なし。再生成のたびに新しい Key で保存するので、古い音声は配信されない | [architecture.md A-18](../architecture.md#a-18) |
