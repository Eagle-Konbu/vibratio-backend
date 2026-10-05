# MVP の範囲外とした拡張

将来必要になる可能性があるが、MVP には入れないもの。追加するときに影響を受ける概念と、関連する判断を示す。

| 拡張 | 必要になる状況 | 影響を受ける概念 | 関連する判断 |
| --- | --- | --- | --- |
| Article の蓄積 | 過去記事の検索、既読管理、記事へのフィードバック | Article が Entity になり、Repository と重複判定が必要になる | [D-01](article.md#d-01) |
| 過去の Episode との URL 照合 | 期間で絞っても日をまたいだ重複が気になる | EpisodeRepository に URL 検索を追加する | [D-04](article.md#d-04) |
| URL の正規化 | 完全一致では重複を除けない例が出てくる | SelectArticles の重複除去 | [D-05](article.md#d-05) |
| 記事のインデックスを LLM に出力させる | 不正な Reference が頻発する | ScriptGenerator の出力、Use Case の検証 | [D-03](generate-episode.md#d-03) |
| 複数の Source に載った記事を重要度のシグナルにする | 記事選定の質を上げたい | Article に掲載 Source 数を持たせる | [D-05](article.md#d-05) |
| Source ごとに異なる上限 k | 特定の Source の量を調整したい | Source に属性を追加する | ― |
| Source の Enabled | CMS から一時的に Source を止めたい | Source に属性を追加する | ― |
| フィード以外の Source（X、GitHub など） | フィードを提供しない情報源を使いたい | Source に種類と種類ごとの設定を導入する | [D-09](source.md#d-09) |
| Topic と本文の位置の対応（チャプター） | 話題ごとに再生位置を移動したい | Topic に位置またはチャプター時刻を追加する | [D-02](episode.md#d-02) |
| Audio を後から付ける | Script を手で修正してから音声化したい | Audio を任意にする | [D-13](episode.md#d-13) |
| Duration | 一覧での表示、Podcast 配信 | Audio に属性を追加する | [D-12](episode.md#d-12) |
| 過去の日付を指定した再生成 | CMS から任意の日の Episode を作り直したい | 基準時刻を引数に取る Use Case | [D-15](episode.md#d-15) |
| CMS での Episode の閲覧 | 過去の Episode を Discord 以外から見たい | GraphQL の Query。Episode は DynamoDB に保存するため、これを見越した構成になっている | ― |
