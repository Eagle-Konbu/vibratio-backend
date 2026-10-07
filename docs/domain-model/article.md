# Article と記事の選定

## Article

Source から取得した記事を表す Value Object。**生成処理の中だけで使う入力であり、保存しない。**

```text
Article
- SourceName   取得元 Source の Name（値のコピー）
- Title
- URL
- PublishedAt  必須
- Summary      プレーンテキスト。選定後は N 文字以内。空でもよい
```

### Identity を持たない

Article は日をまたいで蓄積・管理する対象ではない（[D-01](#d-01)）。そのため ID も Repository も持たない。

<a id="d-01"></a>

> **D-01: Article は保存しない**
>
> - 検討した案: Article を保存する Entity にする / 取得結果を Infrastructure にスナップショットとして保存する
> - 理由: 1日1回のバッチで、入力はフィードのメタデータだけである。その日の取得結果で処理は完結し、蓄積しなくても機能は成り立つ。入口としての価値（どの記事を紹介したか）は、Episode が URL と Title を値として持てば満たせる
> - 見直す条件: 過去記事の検索、既読管理、記事へのフィードバックが必要になったとき

### SourceName

番組の中で「〇〇によると」と情報源に触れるために、LLM への入力として持つ。Source の ID を持たずに Name の値をコピーして持つので、保存しない Article が Source Aggregate に依存せずに済む（[D-08](#d-08)）。

<a id="d-08"></a>

> **D-08: Article は SourceName を値として持つ**
>
> - 検討した案: Source の情報を持たない / Reference にも SourceName を残す
> - 理由: SourceName があれば番組の中で情報源に触れられ、ラジオとしての質が上がる。値のコピーなので D-01 と矛盾しない。Reference の出どころは URL のドメインから分かる

### PublishedAt

必須。選定ルールが公開日時で対象期間を判定するため、PublishedAt を持たない記事は Domain に入らない。フィードから公開日時を得られない記事の扱い（除外する、Atom の `updated` で代用する等）は FeedFetcher の Adapter が決める（[D-06](#d-06)）。

<a id="d-06"></a>

> **D-06: PublishedAt は必須**
>
> - 検討した案: 取得時刻で代用する / 任意項目にして期間の判定から外す
> - 理由: D-04 の対象期間を判定できない記事は扱えない。取得時刻で代用すると毎日同じ記事が新着として入り、D-04 の前提が崩れる。フィード形式ごとの代替（Atom の `updated` など）は Adapter が判断する

### Summary

フィードから得られる**概要**であり、**プレーンテキストで、選定後は N 文字以内**と定義する（[D-07](#d-07)）。

- この定義により、LLM に渡す記事が本文を持たないことと、記事1件あたりの Summary の量に上限があることを保証する
- どの要素から取るか（RSS `description`、Atom `summary` / `content`）と HTML の除去は FeedFetcher の Adapter の責務
- N 文字への切り詰めは SelectArticles が行う（[記事の選定](#記事の選定selectarticles)）。Source の種類や Adapter の実装によらず、上限を1か所で保証するため
- Atom の `content` のように全文が入っている要素も、切り詰めれば冒頭の要約として使える
- 概要が得られない場合は空を許容する。Web ページを取得して要約を生成することはしない

N の具体値は Application の[設定値](generate-episode.md#設定値)とする。

<a id="d-07"></a>

> **D-07: Summary はプレーンテキストで N 文字以内**
>
> - 検討した案: `content` は使わない / 制約を設けない / N 文字への切り詰めを FeedFetcher の Adapter で行う
> - 理由: フィードの `description` や `content` には全文や HTML が入っていることが多い。この定義によって「本文を持たない」と「記事1件あたりの Summary の量に上限がある」を構造的に保証する。Title などのメタデータは制限しないため、入力量全体の上限にはならない。冒頭を切り詰めるだけなら本文を持つことにはならない。切り詰めを Adapter に任せると、LLM のコストを抑える保証が Adapter の正しさに依存し、Source の種類を増やしたときに漏れうる。SelectArticles で行えば、上限は Domain の1か所で保証され、選ばれなかった記事を切り詰める無駄もない

### 本文を持たない理由

[D-01](#d-01)、[D-07](#d-07) の前提でもある。

- LLM への入力コストを抑える
- フィードの取得だけで処理が完結する
- ラジオは一次情報の代替ではなく、詳細はユーザーが URL から原記事を読む

## 記事の選定（SelectArticles）

取得した Article から、その日の番組の素材を決める Domain の純粋関数。

```text
SelectArticles(fetched: Source ごとの Article[], now, window, k, n) → Article[]
```

ルールは次の順に適用する。

1. **対象期間**: `now - window` 以降に公開された記事だけを残す
2. **Source ごとの上限**: Source ごとに PublishedAt の新しい順で k 件までに絞る
3. **重複除去**: URL が完全一致する記事は1件にまとめる（先に現れた Source の記事を残す）
4. **Summary の切り詰め**: 残った記事の Summary を n 文字（rune 単位）までに切り詰める（[D-07](#d-07)）

### 設計上のポイント

- 対象期間で絞るので、前日に紹介した記事が再び候補になること（日をまたいだ重複）はほぼない。そのため過去の Episode との URL 照合は行わない（[D-04](#d-04)）
- RSS や LLM の仕様に依存しないルールなので Domain に置き、Port のモックなしで単体テストする（[D-05](#d-05)）
- 入力を Source ごとにまとめて受け取るので、SourceName が重複していても Source ごとの上限を正しく適用できる
- window・k・N の具体値は Application の[設定値](generate-episode.md#設定値)として引数で渡す
- 重複除去は Source の順序に依存する。順序を決定的にするのは呼び出す側（Use Case）の責務である
- LLM への入力量は「Source 数 × k × (SourceName + Title + URL + PublishedAt + N)」で見積もれる。k と N が制限するのは記事数と Summary の長さだけであり、入力全体に厳密な上限はない。SourceName・Title・URL・PublishedAt はフィードのメタデータであり、Domain では長さを制限しない。Source 数は CMS で登録した数で決まる

<a id="d-04"></a>

> **D-04: 対象期間は直近 24 時間**
>
> - 検討した案: 前回の生成以降を対象にする / 過去の Episode の References と照合する / 期間で絞らない
> - 理由: 1日1回の実行なら、期間で絞るだけで日をまたいだ重複はほぼ起きない。「前回の生成以降」は「前回の生成日時」という状態を必要とする。URL 照合は Repository に検索機能を求め、MVP の範囲を広げる
> - 見直す条件: 期間で絞っても日をまたいだ重複が気になるとき（[future.md](future.md)）

<a id="d-05"></a>

> **D-05: 記事の選定ルールは Domain の関数**
>
> - 検討した案: Use Case の手続きとして書く / 不変条件を持つ Collection 型を作る
> - 理由: これらは RSS や LLM の仕様に依存しない「番組の素材をどう選ぶか」のルールである。外部サービスとの境界の問題である D-03 とは性質が異なる。関数にすればルールが1か所にまとまり、モックなしでテストできる。Collection 型は関数に比べて得るものが少ない
> - 見直す条件: URL の完全一致では重複を除けない例が出てきたとき、複数の Source に載った記事を重要度に使いたくなったとき（[future.md](future.md)）
