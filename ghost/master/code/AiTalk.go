package main


import (
    "fmt"
    "regexp"
    "strings"
    "context"
    "os"
    "os/signal"
    //"math/rand"
    "time"
    openai "github.com/sashabaranov/go-openai"
)


//ChatHistoryArray = append( ChatHistoryArray, openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleSystem   , Content: "" } )
//ChatHistoryArray = append( ChatHistoryArray, openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleUser     , Content: "" } )
//ChatHistoryArray = append( ChatHistoryArray, openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleAssistant, Content: "" } )
//msgs = append( msgs, openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleAssistant, Content: resp.Choices[0].Message.Content } )

type ChatGPTConfig struct {
    API_KEY         string
    //API初期値
    ChargeAPI       int 
    //一分間あたりの使用可能数
    ChargeAPIMax    int
    //補充クールタイム。
    ChargeAPISecMax int
    ChatHistoryMax  int 
    NoticeText      string
}
var MyChatGPTConfig ChatGPTConfig


var ThreadUse     int = 0

//地の文を減らしたい。
var REPLACE_TEXT    = regexp.MustCompile("\\(.*?\\)")


//使用する履歴数を確保する。
var ChatHistoryArray []openai.ChatCompletionMessage



//IDで履歴を分けても良いかもしれないが、
//別の対話で気になったモノを質問するのはめんどいかな。
//当分共通化しておくか。
func AiTalk( NewMessage string, ID string ) {
    if NewMessage == "" {
        return
    }


    fmt.Println( "-------------" )
    fmt.Println( "Talk AI :" + ID )
    if ThreadUse == 1 {
        fmt.Println( "Thread SKIP" )
        NextSecondTalk = append( NextSecondTalk , "Thread使用中 :" + NewMessage )
        return 
    }
    if MyChatGPTConfig.ChargeAPI <= 0 {
        fmt.Println( "API STOP" )
        NextSecondTalk = append( NextSecondTalk , "API枯渇 :" + NewMessage )
        return  
    }
    MyChatGPTConfig.ChargeAPI = MyChatGPTConfig.ChargeAPI - 1
    ThreadUse = 1

    ////client.go
    client := openai.NewClient( MyChatGPTConfig.API_KEY )
    ////chat.go
    msgs := []openai.ChatCompletionMessage{}
    //性格とルールの注入
    msgs = SetPersonYuzuki( msgs )

    //溢れている分を古いものから削除
    for( len( ChatHistoryArray ) > MyChatGPTConfig.ChatHistoryMax ) {
        ChatHistoryArray   = ChatHistoryArray[1:]
    }

    //会話履歴を注入する。
    n := 0
    for ( len( ChatHistoryArray ) != n ){
        //古いものから挿入する。
        msgs = append( msgs, ChatHistoryArray[ n ] )
        n++
    }

    //今回のセリフ。
    msgs = append( msgs , openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleUser, Content: NewMessage })

    //条件等はここで指定できるな。
    Sig, cancel := signal.NotifyContext( context.Background() , os.Interrupt )
    defer cancel()
    resp, err := client.CreateChatCompletion(
        Sig,
        openai.ChatCompletionRequest{ 
            Model   : openai.GPT3Dot5Turbo, 
            Messages: msgs  ,
            //MaxTokens        : 600  ,
            //元->設定なし。
            //Temperature      : 0.1,
            //まともだけど、振れ幅がない。
            //Temperature      : 0.5,
            //0.6からお仕事お疲れさまでしたか?という謎文章になる。
            Temperature      : 0.6,


            //Temperature      : 1.1,
            //PresencePenalty  : -0.6,
            //Temperature      : 1.1,
            //PresencePenalty  : -1.0,
        },)


    if err != nil {
        fmt.Println( "ChatGPT関係のエラー" )
        //rate limit抵触で発生した。
        //トークンの消費数超過でも発生した。4097トークン。
        fmt.Fprintln(os.Stderr, err)
        ThreadUse = 0
        return 
        //return "" , errors.New( "ChatGPT_Error" )
    }



    TextChatGPT := resp.Choices[0].Message.Content
    //%他、危険そうなものは置換しておく。
    TextChatGPT  = strings.Replace( TextChatGPT , "%" , "パーセント" , -1 )
    TextChatGPT = strings.Replace( TextChatGPT , "\\\\" , "\\" , -1 )
    //TextChatGPT  = strings.Replace( TextChatGPT , "" , "" , -1 )


    Check       := regexp.MustCompile( `\\s\[.*?]` )
    //Check       := regexp.MustCompile( `\\s\[\d+?\]` )
    Match       := Check.FindAllStringSubmatch( TextChatGPT , -1 )
    //マッチ箇所があれば。
    //表情差分は必ず出力してくれるわけではない。
    if len( Match ) != 0 {
        //マッチ部分を%sに置き換え
        TextChatGPT = Check.ReplaceAllString( TextChatGPT , "%s" )

        //%sから始まらないなら追加
        StartCheck       := regexp.MustCompile( `^%s` )
        StartBool   := true
        if !StartCheck.MatchString( TextChatGPT ) {
            StartBool = false
            TextChatGPT = "%s" + TextChatGPT
        }
        //先頭に追加。

        //sprintfの受け付ける形。
        ListSurfaceChatGPT := []interface{}{}
        //flat化してリストに。
        for _, v := range Match {
            ListSurfaceChatGPT = append( ListSurfaceChatGPT , v[0] )
        }

        SurfaceTop := []interface{}{}
        //追加したなら数合わせ
        if !StartBool {
            SurfaceTop = append( SurfaceTop , ListSurfaceChatGPT[0] )
        }
        ListSurfaceChatGPT = append( SurfaceTop , ListSurfaceChatGPT[0:]... )
        TextChatGPT = fmt.Sprintf( TextChatGPT , ListSurfaceChatGPT... )
    } else {
    }
    //fmt.Println( TextChatGPT )
    //実際に発言するフェイズ
    //加工して保存する。
    //TextChatGPT  = REPLACE_TEXT.ReplaceAllString( TextChatGPT , "" )
    ChatHistoryArray = append( ChatHistoryArray , openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleUser, Content: NewMessage })
    ChatHistoryArray = append( ChatHistoryArray , openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleAssistant, Content: TextChatGPT } )

    //検閲配列
    CheckArray := []string{ "AI" , "人工知能" }
    for _,v := range CheckArray {
        //AIという単語が入っていた場合履歴に残さない。
        var CheckAI = regexp.MustCompile( v )
        if CheckAI.MatchString( TextChatGPT )  {
            //送信と受信で二つ分。
            ChatHistoryArray   = ChatHistoryArray[:len( ChatHistoryArray ) - 2 ]
            fmt.Println( "AIワードを削除。" )
            break
        }
    }


    TextChatGPT = "\\0\\b[2]" + TextChatGPT
    TextChatGPT = strings.Replace( TextChatGPT , "\\s" , "\\w9\\w3\\s" , -1 )
    SaveTalk( "" )
    SaveTalk( NewMessage )
    SaveTalk( TextChatGPT )

    NextSecondTalk = append( NextSecondTalk , TextChatGPT )
    //fmt.Println( "ここまでのログ" )
    //fmt.Println( ChatHistoryArray )
    ThreadUse = 0
}



//もっとわかりやすい会話の例文を入れて自分好みに染める。
func SetPersonYuzuki( msgs []openai.ChatCompletionMessage ) []openai.ChatCompletionMessage {
    NOW  := time.Now().Format( "PM 2006/01/02 15:04" )

    msgs = append( msgs, openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleSystem, Content: `
    ` + DefaultPrompt + `
    今の時間は` + NOW + `です。時間に合わせた挨拶や言動をしなさい。
    直接現在時刻を言う必要はありません。
    ` } )

    return msgs
}


//double clickに積んでいる。
func AiDelete(){
	if len( ChatHistoryArray ) >= 2 {
		ChatHistoryArray   = ChatHistoryArray[:len( ChatHistoryArray ) - 2 ]
	}
}


func TopicTalk( ID string ) {
    //質問に関しては最初のプロンプトにルールを指定する方が、履歴に残らなくてクリーンだろう。
    //感情表現のルールとかを残すかどうかが焦点かな。

    //前の会話が無ければ終了
    if 2 > len( ChatHistoryArray ) {
        return
    }

    fmt.Println( "-------------" )
    fmt.Println( "Talk AI :" + ID )
    if ThreadUse == 1 {
        fmt.Println( "Thread SKIP" )
        NextSecondTalk = append( NextSecondTalk , "Thread使用中 :" + ID )
        return 
    }
    if MyChatGPTConfig.ChargeAPI <= 0 {
        fmt.Println( "API STOP" )
        NextSecondTalk = append( NextSecondTalk , "API枯渇 :" + ID )
        return  
    }
    MyChatGPTConfig.ChargeAPI = MyChatGPTConfig.ChargeAPI - 1
    ThreadUse = 1

    client := openai.NewClient( MyChatGPTConfig.API_KEY )
    msgs := []openai.ChatCompletionMessage{}

    //会話履歴を注入する。
    msgs = append( msgs, ChatHistoryArray[ len( ChatHistoryArray ) - 1] )
    //性格とルールの注入
    msgs = SetPersonTopic( msgs )
    fmt.Println( " log " )
    fmt.Println( ChatHistoryArray[ len( ChatHistoryArray ) - 1] )


    //今回のセリフ。
    //msgs = append( msgs , openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleUser, Content: NewMessage })

    //条件等はここで指定できるな。
    Sig, cancel := signal.NotifyContext( context.Background() , os.Interrupt )
    defer cancel()
    resp, err := client.CreateChatCompletion(
        Sig,
        openai.ChatCompletionRequest{ 
            Model   : openai.GPT3Dot5Turbo, 
            Messages: msgs  ,
            //MaxTokens        : 600  ,
            //元->設定なし。
            //Temperature      : 0.1,
            //まともだけど、振れ幅がない。
            //Temperature      : 0.5,
            //0.6からお仕事お疲れさまでしたか?という謎文章になる。
            Temperature      : 0.6,
            //Temperature      : 1.1,
            //PresencePenalty  : -0.6,
            //Temperature      : 1.1,
            //PresencePenalty  : -1.0,
        },)


    if err != nil {
        fmt.Println( "ChatGPT関係のエラー" )
        //rate limit抵触で発生した。
        //トークンの消費数超過でも発生した。4097トークン。
        fmt.Fprintln(os.Stderr, err)
        ThreadUse = 0
        return 
        //return "" , errors.New( "ChatGPT_Error" )
    }



    TextChatGPT := resp.Choices[0].Message.Content
    //%他、危険そうなものは置換しておく。
    TextChatGPT  = strings.Replace( TextChatGPT , "%" , "パーセント" , -1 )
    TextChatGPT = strings.Replace( TextChatGPT , "\\\\" , "\\" , -1 )
    //TextChatGPT  = strings.Replace( TextChatGPT , "" , "" , -1 )


    //会話の履歴から、表情が発生しないとは言い切れない。
    Check       := regexp.MustCompile( `\\s\[.*?]` )
    //Check       := regexp.MustCompile( `\\s\[\d+?\]` )
    Match       := Check.FindAllStringSubmatch( TextChatGPT , -1 )
    //マッチ箇所があれば。
    //表情差分は必ず出力してくれるわけではない。
    if len( Match ) != 0 {
        //マッチ部分を%sに置き換え
        TextChatGPT = Check.ReplaceAllString( TextChatGPT , "%s" )

        //%sから始まらないなら追加
        StartCheck       := regexp.MustCompile( `^%s` )
        StartBool   := true
        if !StartCheck.MatchString( TextChatGPT ) {
            StartBool = false
            TextChatGPT = "%s" + TextChatGPT
        }
        //先頭に追加。

        //sprintfの受け付ける形。
        ListSurfaceChatGPT := []interface{}{}
        //flat化してリストに。
        for _, v := range Match {
            ListSurfaceChatGPT = append( ListSurfaceChatGPT , v[0] )
        }

        SurfaceTop := []interface{}{}
        //追加したなら数合わせ
        if !StartBool {
            SurfaceTop = append( SurfaceTop , ListSurfaceChatGPT[0] )
        }
        ListSurfaceChatGPT = append( SurfaceTop , ListSurfaceChatGPT[0:]... )
        TextChatGPT = fmt.Sprintf( TextChatGPT , ListSurfaceChatGPT... )
    } else {
    }
    //OnCommunicate等で直近の会話の整合性を保つ。
    ChatHistoryArray[ len(ChatHistoryArray) - 1 ] = openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleAssistant, Content: TextChatGPT } 
    //ChatHistoryArray = append( ChatHistoryArray , openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleAssistant, Content: TextChatGPT } )
    
    //話題確保
    if regexp.MustCompile(`【.*?】`).MatchString( TextChatGPT ) {
        word := regexp.MustCompile(`^.*?【`).ReplaceAllString( TextChatGPT , "" )
        word  = regexp.MustCompile(`】.*`).ReplaceAllString( word , "" )

        text := strings.Split( word , "\n" )
        if len( text ) != 1 {
            word  = text[0]
            fmt.Println( "-----------wordlist-------" )
            WordList = append( WordList , word )
            fmt.Println( word )
            fmt.Println( WordList )
            fmt.Println( "-----------wordlist-------" )
        }
    }


    TextChatGPT = "\\0\\b[2]" + TextChatGPT
    TextChatGPT = strings.Replace( TextChatGPT , "\\s" , "\\w9\\w3\\s" , -1 )
    SaveTalk( "" )
    SaveTalk( TextChatGPT )

    NextSecondTalk = append( NextSecondTalk , TextChatGPT )
    //fmt.Println( "ここまでのログ" )
    //fmt.Println( ChatHistoryArray )
    ThreadUse = 0
}


//この時点で履歴を詰め込んでいないと適当にしゃべり始めた。
func SetPersonTopic( msgs []openai.ChatCompletionMessage ) []openai.ChatCompletionMessage {
    i := 0
    msg := "会話の履歴から、適当な固有名詞や単語を選択し、その背景や内容、類似作品など話を掘り下げなさい。"
    if len( WordList ) != 0 {
        msg = msg + "ただし、下記のリストの単語は選択してはいけません。["
    }
    for( i != len( WordList ) ){
        msg = msg + WordList[ i ] + ","
        i++
    }
    if len( WordList ) != 0 {
        msg = msg + "]"
    }
    x := `
    選択した単語は会話の一行目に配置し、【】で囲い改行しなさい。例えば、【vim】についてお話ししましょう。といった風に始め、2行目からその単語の説明をしてください。
    説明文は、尻切れや中途半端になってはいけません。
    `
    msg = msg + x
    msgs = append( msgs, openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleSystem, Content: msg })
    return msgs
}







