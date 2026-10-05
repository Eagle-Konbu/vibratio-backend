# Article と記事の選定

## Article

Source から取得した記事を表す Value Object。**生成処理の中だけで使う入力であり、保存しない。**

```text
Article
- SourceName   取得元 Source の Name（値のコピー）
- Title
- URL
- PublishedAt  必須
- Summary      プレーンテキスト、N 文字以内。空でもよい
```

### Identity を持たない

Article は日をまたいで蓄積・管理する対象ではない（[D-01](decisions.md#d-01-article-は保存しない)）。そのため ID も Repository も持たない。

### SourceName

番組の中で「〇〇によると」と情報源に触れるために、LLM への入力として持つ。Source の ID を持たずに Name の値をコピーして持つので、保存しない Article が Source Aggregate に依存せずに済む（[D-08](decisions.md#d-08-article-は-sourcename-を値として持つ)）。

### PublishedAt

必須。選定ルールが公開日時で対象期間を判定するため、PublishedAt を持たない記事は Domain に入らない。フィードから公開日時を得られない記事の扱い（除外する、Atom の `updated` で代用する等）は FeedFetcher の Adapter が決める（[D-06](decisions.md#d-06-publishedat-は必須)）。

### Summary

フィードから得られる**概要**であり、**プレーンテキストで N 文字以内**と定義する（[D-07](decisions.md#d-07-summary-はプレーンテキストで-n-文字以内)）。

- この定義により、記事本文を持たないことと、記事1件あたりの Summary の量に上限があることを保証する
- どの要素から取るか（RSS `description`、Atom `summary` / `content`）、HTML の除去、切り詰めは FeedFetcher の Adapter の責務
- Atom の `content` のように全文が入っている要素も、切り詰めれば冒頭の要約として使える
- 概要が得られない場合は空を許容する。Web ページを取得して要約を生成することはしない

N の具体値は Application の[設定値](generate-episode.md#設定値)とする。

### 本文を持たない理由

[D-01](decisions.md#d-01-article-は保存しない)、[D-07](decisions.md#d-07-summary-はプレーンテキストで-n-文字以内) の前提でもある。

- LLM への入力コストを抑える
- フィードの取得だけで処理が完結する
- ラジオは一次情報の代替ではなく、詳細はユーザーが URL から原記事を読む

## 記事の選定（SelectArticles）

取得した Article から、その日の番組の素材を決める Domain の純粋関数。

```text
SelectArticles(fetched: Source ごとの Article[], now, window, k) → Article[]
```

ルールは次の順に適用する。

1. **対象期間**: `now - window` 以降に公開された記事だけを残す
2. **Source ごとの上限**: Source ごとに PublishedAt の新しい順で k 件までに絞る
3. **重複除去**: URL が完全一致する記事は1件にまとめる（先に現れた Source の記事を残す）

### 設計上のポイント

- 対象期間で絞るので、前日に紹介した記事が再び候補になること（日をまたいだ重複）はほぼない。そのため過去の Episode との URL 照合は行わない（[D-04](decisions.md#d-04-対象期間は直近-24-時間)）
- RSS や LLM の仕様に依存しないルールなので Domain に置き、Port のモックなしで単体テストする（[D-05](decisions.md#d-05-記事の選定ルールは-domain-の関数)）
- 入力を Source ごとにまとめて受け取るので、SourceName が重複していても Source ごとの上限を正しく適用できる
- window・k の具体値は Application の[設定値](generate-episode.md#設定値)として引数で渡す
- LLM への入力量は「Source 数 × k × (SourceName + Title + URL + PublishedAt + N)」で見積もれる。このうち上限を定めるのは k と N（Summary）だけである。SourceName・Title・URL・PublishedAt はフィードのメタデータで、Domain では長さを制限しない。Source 数も CMS で登録した数で決まり、上限を設けない
