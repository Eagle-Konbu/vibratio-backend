# Source

記事の情報源を表す Entity。CMS で管理する。

```text
Source
- ID       CMS が採番する識別子
- Name     表示名。LLM への入力（SourceName）にも使う
- FeedURL  RSS / Atom フィードの URL
```

## ライフサイクル

- CMS（AppSync 経由）で作成・更新・削除する。backend のバッチは読み取りのみ
- Episode は Source を参照しないため、Source を削除しても過去の Episode に影響しない

## 持たないもの

| 属性 | 持たない理由 |
| --- | --- |
| Type / FeedFormat（RSS / Atom） | フィード形式は取得した XML から自動判別できる。人が選ぶ項目を残すと、実際の形式と食い違うデータが入り込む余地が生まれる（[D-10](decisions.md#d-10-source-はフィード形式を持たない)） |
| Endpoint（汎用的な名前） | MVP の Source はフィードだけであり、意味の明確な FeedURL とする。X など取得方式が異なる Source は、追加する時点で構造を設計する（[D-09](decisions.md#d-09-source-の種類を抽象化しない)） |
| Enabled | MVP では不要。持たせるかどうかは CMS の実装時に決める（未決。[future.md](future.md)） |

## 永続化

Source の永続化は Infrastructure の責務である。現在は DynamoDB の config テーブルに保存している。属性名 `endpoint` は Repository で FeedURL に対応づける。
