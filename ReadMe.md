# ukagakaAiGhost_Yukari
このゴーストは、誰でもAI搭載ゴーストを作れるように作成されました。
しかし、あくまで受動的なゴーストです。


## うちの子をAI化したいのだけど?
## サンプル動画とかある?


## 設定ファイルの場所
`/ghost/master/Config.txt`


#### 設定項目
```
{
     "API_KEY"          : "",
     "NoticeText"       : "\\1!\\0" ,
     "ChargeAPI"        : 2  , 
     "ChargeAPIMax"     : 3  ,
     "ChargeAPISecMax"  : 30 ,
     "ChatHistoryMax"   : 7 
}
```

- API_KEY          
    ChatGPT API Keyを入力してください。
- NoticeText       
    ユーザのアクション -> 通信 -> 通信内容の表示となるため通信が始まったことを通知するテキストを設定出来るようにしています。
- ChargeAPI        
    起動時のチケット枚数です。
    ChatGPT APIは、一分間に3回までしか使用できないのでこのように管理します。
    通信するたびに一枚消費します。足りないときは通信されず、ゴーストに通知されます。
- ChargeAPIMax     
    チケットの最大値です。
    時間経過で補充されます。
- ChargeAPISecMax  
    何秒おきで補充するか設定する項目です。
- ChatHistoryMax   
    ChatGPTとの会話履歴を保存する数です。
    数値を1増やすごとに質問と回答のセットを保存し、通信に追加します。
    増やすとコミュニケーションの正確性が増しますが、通信量が増え、トークンの消費が増えます。
    一度の通信によるトークン数には上限があるため7程度が良いでしょう。


## プロンプト
デフォルトプロンプトと各種触り判定用のプロンプトがあります。
デフォルトプロンプトは毎回、通信の最初に使用されるもので、全文読み込まれます。
ChatGPTには回答や喋り方に振れ幅がある(そのように設定している。)ので返答の例文を入れておくと安定します。

各種触り判定のプロンプトは対応するテキストファイルからランダムに一行を使用します。
改行のみの行はスキップされます。

#### 例えば
頭を撫でまわした場合、下記のようになります。

1. デフォルトプロンプト
1. 現在の時刻
1. 会話履歴 x N
1. 撫で反応プロンプト(ランダムに1行)

これらを順番につなげたものがChatGPTに送信されます。


#### デフォルトプロンプトの場所
`/ghost/master/DefaultPrompt.txt`

#### 実装されている各種触り判定と使用されるプロンプトの場所
- DoubleClick         頭
- DoubleClick         胸
- 撫でまわし          頭
- 撫でまわし          胸
- ホイールクリック    口

```
├─ghost
│  └─master
│      ├─DoubleClick
│      │      BustPrompt.txt
│      │      HeadPrompt.txt
│      ├─Nade
│      │      BustPrompt.txt
│      │      HeadPrompt.txt
│      └─WheelClick
│              MousePrompt.txt
```


#### 特殊な実装がされている判定
- 設定部位がないところでホイール下
    最後の会話から話題を選んで掘り下げる。
    その後、一分間隔でさらに掘り下げ続ける。
    一度選んだ話題は使われなくなる。
    特性上、一番最初に実行する際は、前回の会話と同じような内容になってしまう。
    XXXXってなに?と聞いた後に使うと延々と雑学を垂れ流してくれる。

- 設定部位がないところでホイール上
    会話の掘り下げを終了する。

- 設定部位がないところでDoubleClick
    最後の会話履歴を削除する。
    AIが妙な事を言いだし、それが履歴に残ると引きずり続けた会話をするため、手動で消せるようにした。


#### 特殊な挙動
「AI」および「人工知能」が含まれた会話は自動で削除し、履歴に残らないようにしている。
私はAIなので～～が出来ない。と言い始めるとやはり引きずるので自動化している。


## ログ
送信した内容と受信した内容を下記に保存しています。
作りたいゴーストのプロンプトを入力しておくことで、既存のゴーストの作成効率を上げることができるでしょう。
`/ghost/master/log/log.txt`


## 同梱されている立ち絵について
こちらの立ち絵をお借りしました。
お名前がうまく打てない。申し訳ない。
[【立ち絵フリー素材】VOICEROID結月ゆかり 中古0円 - ニコニコ静画 (イラスト)](https://seiga.nicovideo.jp/seiga/im11102824)


## 問題点
golangで書かれたdllの為、freelibrary等をするとSSPが落ちます。
具体的には、ゴーストを終了した際に、他のゴーストを残した状態でもSSPが落ちてしまいます。

プロンプトファイルが見つからなかった場合、エラーで落ちます。
通信環境がなかった場合、エラーで落ちます。


## 注意書き
このプログラムなどを使用したことによるいかなる問題や損害に対して私は責任を負いません。

                
## Author
ambergon
