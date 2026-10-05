# GenerateEpisode Use Case

Domain Model から導かれる Application 層の構成。MVP の Use Case はこの1つだけである。

## トリガー

- 毎朝のスケジュール実行
- 失敗時や再生成時の手動実行（同じ Use Case を使う）

## 手順

```text
1. SourceRepository.FindAll()                          → Source[]
2. FeedFetcher.Fetch(source)  ※Source ごと             → Article[]
3. SelectArticles(fetched, now, window, k)            → Article[]
     0件なら Episode を作らずに正常終了する（D-16）
4. ScriptGenerator.Generate(articles)                  → Script, Topics
5. Reference を入力記事と照合する
     入力記事に存在しない URL の Reference を除外し、除外件数をログに出す
     Reference の Title は入力記事の値で上書きする
6. EpisodeRepository.FindByDate(today)
     あれば既存の ID を使い、なければ新しい ID を採番する
7. Synthesizer.Synthesize(script)                       → 音声データ
8. AudioStorage.Put(id, 音声データ)                      → Audio.Key
9. 既存なら Episode.Replace、なければ Episode.New
10. EpisodeRepository.Save(episode)
11. Notifier.Notify(episode)                             → Discord
```

`now` は実行時刻、`today` はその JST の日付である。

## 失敗時

自動では再試行せず、手動で全体を再実行する（[D-13](episode.md#d-13)）。

- 手順10（保存）までに失敗した場合、Episode は保存されない
- 手順11（通知）で失敗した場合、Episode は保存済みである。再実行すると、その日の Episode を再生成して差し替える（Replace）

## LLM 出力の検証

手順5の照合は LLM の出力に対する防御であり、入力記事の集合を持っている Use Case が行う（[D-03](#d-03)）。

不正な Reference があっても Episode の生成は止めない（[D-14](episode.md#d-14)）。頻発する場合は、URL ではなく入力記事のインデックスを LLM に出力させる方式に切り替える（[future.md](future.md)）。

<a id="d-03"></a>

> **D-03: LLM 出力の検証は Use Case の責務**
>
> - 検討した案: Episode の不変条件にする / 検証しない / URL ではなく入力記事のインデックスを LLM に出力させる
> - 理由: この検証は番組そのもののルールではなく、LLM という外部サービスの出力を信用しないための防御である。入力記事の集合は生成処理の中にしか存在しない（D-01）。Episode の不変条件にすると、入力記事の集合を Episode に渡す必要がある。Episode は保存された番組としての整合性だけを保証する。インデックス方式は、不正な Reference が頻発したときの切り替え先として残す
> - 見直す条件: 不正な Reference が頻発したとき（[future.md](future.md)）

## Repository

| Repository | 操作 | 備考 |
| --- | --- | --- |
| SourceRepository | FindAll | 読み取りのみ。書き込みは CMS が行う |
| EpisodeRepository | FindByDate, Save | Save は Date が同じ Episode を上書きする。Date の一意性は Save が原子的に保証する（[episode.md](episode.md#identity-と-date)） |

Article には Repository を置かない。

## Port

| Port | 責務 | Adapter に閉じ込めるもの |
| --- | --- | --- |
| FeedFetcher | Source から Article[] を取得する | RSS / Atom の形式判別、要素の選択、HTML 除去、Summary の切り詰め、PublishedAt のない記事の扱い |
| ScriptGenerator | Article[] から Script と Topics を生成する | LLM のプロンプト、構造化出力のスキーマ、API |
| Synthesizer | Script から音声データを生成する | TTS の API、文字数制限に応じたテキストの分割と結合 |
| AudioStorage | 音声データを保存し Key を返す | 保存先（S3）、Key の作り方（推測困難な要素を含める）、Content-Type |
| Notifier | Episode を通知する | Discord の形式と文字数制限、Key から配信 URL への変換 |

## 設定値

| 値 | 意味 | 目安 |
| --- | --- | --- |
| window | 対象期間 | 24 時間 |
| k | Source ごとの記事数の上限 | 未決。実装時に決める |
| N | Summary の最大文字数 | 未決。300〜500 文字を目安とする |

## 通知内容

Discord には次の内容を送る。ラジオを聴きながらリンクを開けるように、通知だけで入口として完結させる。

```text
音声の URL
Topics[]
  - 見出し
  - 記事のタイトルと URL
```

Episode はタイトルを持たないので、通知にも番組タイトルは含めない。
