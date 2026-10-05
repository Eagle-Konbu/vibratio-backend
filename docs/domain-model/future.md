# MVP の範囲外とした拡張

将来必要になる可能性があるが、MVP には入れないもの。追加するときに影響を受ける概念と、関連する判断を示す。

| 拡張 | 必要になる状況 | 影響を受ける概念 | 関連する判断 |
| --- | --- | --- | --- |
| Article の蓄積 | 過去記事の検索、既読管理、記事へのフィードバック | Article が Entity になり、Repository と重複判定が必要になる | [D-01](decisions.md#d-01-article-は保存しない) |
| 過去の Episode との URL 照合 | 期間で絞っても日をまたいだ重複が気になる | EpisodeRepository に URL 検索を追加する | [D-04](decisions.md#d-04-対象期間は直近-24-時間) |
| URL の正規化 | 完全一致では重複を除けない例が出てくる | SelectArticles の重複除去 | [D-05](decisions.md#d-05-記事の選定ルールは-domain-の関数) |
| 記事のインデックスを LLM に出力させる | 不正な Reference が頻発する | ScriptGenerator の出力、Use Case の検証 | [D-03](decisions.md#d-03-llm-出力の検証は-use-case-の責務) |
| 複数の Source に載った記事を重要度のシグナルにする | 記事選定の質を上げたい | Article に掲載 Source 数を持たせる | [D-05](decisions.md#d-05-記事の選定ルールは-domain-の関数) |
| Source ごとに異なる上限 k | 特定の Source の量を調整したい | Source に属性を追加する | ― |
| Source の Enabled | CMS から一時的に Source を止めたい | Source に属性を追加する | ― |
| フィード以外の Source（X、GitHub など） | フィードを提供しない情報源を使いたい | Source に種類と種類ごとの設定を導入する | [D-09](decisions.md#d-09-source-の種類を抽象化しない) |
| Topic と本文の位置の対応（チャプター） | 話題ごとに再生位置を移動したい | Topic に位置またはチャプター時刻を追加する | [D-02](decisions.md#d-02-script-と-topics-を分離する) |
| Audio を後から付ける | Script を手で修正してから音声化したい | Audio を任意にする | [D-13](decisions.md#d-13-episode-は完成したときだけ存在する) |
| Duration | 一覧での表示、Podcast 配信 | Audio に属性を追加する | [D-12](decisions.md#d-12-audio-は保存先に依存しない-key-だけを持つ) |
| 過去の日付を指定した再生成 | CMS から任意の日の Episode を作り直したい | 基準時刻を引数に取る Use Case | [D-15](decisions.md#d-15-再生成は同じ-use-case-で上書きする) |
| CMS での Episode の閲覧 | 過去の Episode を Discord 以外から見たい | GraphQL の Query。Episode は DynamoDB に保存するため、これを見越した構成になっている | ― |
