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

- Identity は ID である。外部に公開する識別子（CMS、GraphQL）を日付から切り離すため（[D-11](decisions.md#d-11-episode-の-identity-は-id)）
- Date は生成を実行した JST の日付で、**ユニーク**とする
- 一意性は EpisodeRepository.Save が Date をキーに原子的に保証する。FindByDate で確認してから New する流れだけでは、生成処理が同時に実行されたときに同じ Date の Episode が2件できうるため。同時実行で競合した場合は、後から保存した Episode で上書きする（具体的な実現方法はテーブル設計で決める。[infra-alignment.md](infra-alignment.md)）
- 選定後の記事が0件の日は Episode を作らない（[D-16](decisions.md#d-16-記事が0件の日は-episode-を作らない)）。そのため Episode は1日あたり0件か1件になる
- Episode はタイトルを持たない

## Script と Topics

Script と Topics は、番組の**本文**と**目次**の関係にある（[D-02](decisions.md#d-02-script-と-topics-を分離する)）。

- Script は番組として1つのテキストであり、Topic ごとに分割しない。オープニング・話題間のつなぎ・エンディングなど、どの話題にも属さないテキストを含むため
- Topics は「この回で扱った話題と、その参照記事」を示す。Script のどこで話しているか（位置やチャプター時刻）は持たない
- Topic は ID を持たない。個別に操作する要件がないため

Topics と References は、ユーザーが「ラジオで聞いた話題」から「記事」にたどり着くための構造であり、Discord への通知内容にもなる。

## Reference

Reference は Article の URL と Title を**生成時点でコピーした値**である。Article を参照しないため、Article を保存しなくても Episode 単体で入口として成立する。

References が0件の Topic を許容する。入力記事にない URL を指す Reference（不正な Reference）を除外した結果、リンクのない話題が残ることがある。見出しが残っていれば、検索の手がかりになる（[D-14](decisions.md#d-14-不正な-reference-は除外して続行する)）。

## Audio

Audio は Key だけを持つ Value Object である（[D-12](decisions.md#d-12-audio-は保存先に依存しない-key-だけを持つ)）。

- Key は AudioStorage（Port）が保存時に返す識別子であり、Domain はその中身を解釈しない
- S3 の ObjectKey や配信 URL を Domain に持ち込まない。再生用 URL は Infrastructure が Key から生成する
- Duration は持たない。MVP には取得コストに見合う用途がない

Audio を値として Episode に持たせることで、「存在する Episode は必ず聴ける」を型と不変条件で表す。

## 不変条件

| 条件 | 理由 |
| --- | --- |
| Script が空でない | 番組として成立しない |
| Topics が空でない | 入口として成立しない。紹介する記事がない日は Episode を作らない |
| Audio がある | 存在する Episode は必ず聴ける。Episode に状態（生成中など）を持たせない（[D-13](decisions.md#d-13-episode-は完成したときだけ存在する)） |

LLM の出力が入力記事に由来するかの検証は、Episode の不変条件ではなく Use Case の責務である（[D-03](decisions.md#d-03-llm-出力の検証は-use-case-の責務)）。

## 振る舞い

```text
New(id, date, script, topics, audio) → Episode     不変条件を検証して生成する
Replace(script, topics, audio)                     ID と Date を維持したまま内容を差し替える
```

ID は音声を保存する前（Episode.New / Replace の前）に決める。AudioStorage が ID をもとに Key を作るため、音声を保存する時点で ID が必要になる。

## 再生成

同じ日にバッチを手動で再実行すると、その日の Episode を再生成する（[D-15](decisions.md#d-15-再生成は同じ-use-case-で上書きする)）。

- 既存の Episode があれば、**ID を維持したまま**内容を差し替える
- 再生成前の内容は履歴として残さない
- 再生成でも Discord に通知する
- 過去の日付を指定した再生成は行わない
- 再生成で同じ Key に上書きすると、配信キャッシュに古い音声が残る。対処方法は未決である（[infra-alignment.md](infra-alignment.md)）
