# Episode

1日分の番組を表す Entity であり、Aggregate Root である。

```text
Episode
- ID          Identity
- Date        JST の日付。ユニーク
- Script      番組全体の読み上げテキスト
- Topics[]    話題の目次
- Audio       音声

Topic (VO)
- Title         話題の見出し
- References[]  紹介した記事

Reference (VO)
- URL
- Title

Audio (VO)
- Key         保存先に依存しない識別子
```

## Identity と Date

- Identity は ID である。外部に公開する識別子（CMS、GraphQL）を日付から切り離すため（[D-11](#d-11)）
- Date は生成を実行した JST の日付で、**ユニーク**とする
- 一意性は EpisodeRepository.Save が Date をキーに原子的に保証する。FindByDate で確認してから New するだけでは足りない。生成処理が同時に実行されると、どちらも既存なしと判断し、同じ Date の Episode を2件作りうる
- 同時実行で競合した場合は、後から保存した Episode で上書きする。実現方法はテーブル設計で決める（[infra-alignment.md](infra-alignment.md)）
- 選定後の記事が0件の日は Episode を作らない（[D-16](#d-16)）。そのため Episode は1日あたり0件か1件になる
- Episode はタイトルを持たない

<a id="d-11"></a>

> **D-11: Episode の Identity は ID**
>
> - 検討した案: Date を Identity にする / 同じ Date に複数の Episode を許す
> - 理由: 外部に公開する識別子（CMS、GraphQL、音声の Key）を日付から切り離す。個人利用なので、再生成前の版の履歴は不要である

<a id="d-16"></a>

> **D-16: 記事が0件の日は Episode を作らない**
>
> - 検討した案: 「新しい話題がなかった」という短い Episode を作る / 対象期間を広げて選び直す
> - 理由: 紹介する記事がなければ入口として作る意味がない。短い Episode は TTS のコストがかかり、Topics が空でないという不変条件と矛盾する。期間を広げると前日の記事と重複する

## Script と Topics

Script と Topics は、番組の**本文**と**目次**の関係にある（[D-02](#d-02)）。

- Script は番組として1つのテキストであり、Topic ごとに分割しない。オープニング・話題間のつなぎ・エンディングなど、どの話題にも属さないテキストを含むため
- Topics は「この回で扱った話題と、その参照記事」を示す。Script のどこで話しているか（位置やチャプター時刻）は持たない
- Topic は ID を持たない。個別に操作する要件がないため

Topics と References は、ユーザーが「ラジオで聞いた話題」から「記事」にたどり着くための構造であり、Discord への通知内容にもなる。

<a id="d-02"></a>

> **D-02: Script と Topics を分離する**
>
> - 検討した案: Section ごとに読み上げテキストと References を持つ / References を Episode 単位のフラットな一覧にする / Topic に本文中の位置を持たせる
> - 理由: 番組は1つの作品として保存する。オープニングや話題間のつなぎなど、どの話題にも属さないテキストがあるため、Section 単位に分割するのは不自然である。一方で、話題と記事の対応がないとユーザーが記事を探す手間が残る。位置の対応は LLM の出力からは正確に取りにくく、MVP では使わない
> - 見直す条件: 話題ごとに再生位置を移動したくなったとき（[future.md](future.md)「チャプター」）

## Reference

Reference は Article の URL と Title を**生成時点でコピーした値**である。Article を参照しないため、Article を保存しなくても Episode 単体で入口として成立する。

References が0件の Topic を許容する。入力記事にない URL を指す Reference（不正な Reference）を除外した結果、リンクのない話題が残ることがある。見出しが残っていれば、検索の手がかりになる（[D-14](#d-14)）。

<a id="d-14"></a>

> **D-14: 不正な Reference は除外して続行する**
>
> - 検討した案: 0件になった Topic も除外する / 1件でも不正があれば失敗にする
> - 理由: URL の捏造はたまに起きる程度と想定しており（未計測）、そのたびに毎日のバッチを止めるほどの問題ではない。Script は1つのテキストなので、Topic だけを除外すると読み上げと目次が一致しなくなる。見出しが残れば検索の手がかりになる

## Audio

Audio は Key だけを持つ Value Object である（[D-12](#d-12)）。

- Key は AudioStorage（Port）が保存時に返す識別子であり、Domain はその中身を解釈しない
- S3 の ObjectKey や配信 URL を Domain に持ち込まない。再生用 URL は Infrastructure が Key から生成する
- Duration は持たない。MVP には取得コストに見合う用途がない

Audio を値として Episode に持たせることで、「存在する Episode は必ず聴ける」を型と不変条件で表す。

<a id="d-12"></a>

> **D-12: Audio は保存先に依存しない Key だけを持つ**
>
> - 検討した案: Episode の ID から Infrastructure が規則で導く / 再生用 URL を持つ / Duration を持つ
> - 理由: S3 の ObjectKey や配信 URL は外部サービスの仕様である。それでも Audio を値として持てば、Audio が必須であることを型で表せる。Duration を得るには TTS の結果を解析する必要があり、そのコストに見合う用途（一覧表示、Podcast 配信）は MVP にない
> - 見直す条件: 一覧表示や Podcast 配信で再生時間が必要になったとき（[future.md](future.md)）

## 不変条件

| 条件 | 理由 |
| --- | --- |
| Script が空でない | 番組として成立しない |
| Topics が空でない | 入口として成立しない。紹介する記事がない日は Episode を作らない |
| Audio がある | 存在する Episode は必ず聴ける。Episode に状態（生成中など）を持たせない（[D-13](#d-13)） |

LLM の出力が入力記事に由来するかの検証は、Episode の不変条件ではなく Use Case の責務である（[D-03](generate-episode.md#d-03)）。

<a id="d-13"></a>

> **D-13: Episode は完成したときだけ存在する**
>
> - 検討した案: Script の時点で保存し Audio を後から付ける / 状態（Generating, ScriptReady, ...）を持たせる
> - 理由: 入力がメタデータだけなので、LLM のコストは小さく、全体を再実行しても負担にならないと想定している（未計測）。「存在する Episode は必ず聴ける」が成り立ち、利用側で状態ごとの分岐が要らない
> - 見直す条件: Script を手で修正してから音声化したい、TTS のコストが高く部分的にやり直したいとき

## 振る舞い

```text
New(id, date, script, topics, audio) → Episode     不変条件を検証して生成する
Replace(script, topics, audio)                     ID と Date を維持したまま内容を差し替える
```

ID は音声を保存する前（Episode.New / Replace の前）に決める。AudioStorage が ID をもとに Key を作るため、音声を保存する時点で ID が必要になる。

## 再生成

同じ日にバッチを手動で再実行すると、その日の Episode を再生成する（[D-15](#d-15)）。

- 既存の Episode があれば、**ID を維持したまま**内容を差し替える
- 再生成前の内容は履歴として残さない
- 再生成でも Discord に通知する
- 過去の日付を指定した再生成は行わない
- 再生成で同じ Key に上書きすると、配信キャッシュに古い音声が残る。対処方法は未決である（[infra-alignment.md](infra-alignment.md)）

<a id="d-15"></a>

> **D-15: 再生成は同じ Use Case で上書きする**
>
> - 検討した案: 削除して新しい ID で作り直す / 過去の日付を指定できる再生成専用の Use Case を作る
> - 理由: ID が変わらないので、外部に公開した識別子が壊れない。MVP での再生成は、失敗した日や気に入らなかった日にその場でやり直す程度である。過去の日付は、多くのフィードで記事が消えていて再現できない
> - 見直す条件: CMS から任意の日の Episode を作り直したくなったとき（[future.md](future.md)）
