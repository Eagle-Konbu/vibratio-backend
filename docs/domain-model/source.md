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
| Type / FeedFormat（RSS / Atom） | フィード形式は取得した XML から自動判別できる。人が選ぶ項目を残すと、実際の形式と食い違うデータが入り込む余地が生まれる（[D-10](#d-10)） |
| Endpoint（汎用的な名前） | MVP の Source はフィードだけであり、意味の明確な FeedURL とする。X など取得方式が異なる Source は、追加する時点で構造を設計する（[D-09](#d-09)） |
| Enabled | MVP では不要。持たせるかどうかは CMS の実装時に決める（未決。[future.md](future.md)） |

<a id="d-10"></a>

> **D-10: Source はフィード形式を持たない**
>
> - 検討した案: Domain には持たせず CMS には残す / FeedFormat として Domain に持たせる
> - 理由: フィード形式は取得した XML から自動判別できる。人が選ぶ項目として残すと、実際の形式と食い違うデータを招く。Domain に持たせると XML 形式の知識が Domain に入り込む

<a id="d-09"></a>

> **D-09: Source の種類を抽象化しない**
>
> - 検討した案: 値が RSS だけの Type を持つ / MVP で X も扱う
> - 理由: 値が1つしかない Enum は分岐を生まない。X / GitHub / Reddit などは取得方式がまったく異なり、今 Endpoint で抽象化してもそのまま使える見込みが低い。2種類目を追加する時点で、実際の要件に合わせて設計する。Hacker News や Reddit は当面フィードとして扱える
> - 見直す条件: フィードを提供しない情報源を使いたくなったとき（[future.md](future.md)）

## 永続化

Source の永続化は Infrastructure の責務である。現在は DynamoDB の config テーブルに保存している。属性名 `endpoint` は Repository で FeedURL に対応づける。
