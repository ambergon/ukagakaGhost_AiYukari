package main


import (
    "strconv"
    "math/rand"
)



func OnMouseWheel( numStr string ) string {
    var Value string = ""
    num,_ := strconv.Atoi( numStr )
    if num >= 0 {
        AutoTalkSwitch = false
        Value = "回答モード終了"
    } else {
        AutoTalkSwitch = true
        Value = "回答モード開始"
    }
    return Value
}


func WheelClick( Parts string ) string {
    res := ""
    if Parts == "Mouse" {
        i := rand.Intn( len( WheelClickMousePrompt ) )
        res = WheelClickMousePrompt[ i ]
    }
    if res != "" {
        go AiTalk( res , "OnMouseMove" )
        return MyChatGPTConfig.NoticeText
    }
    return ""
}





func OnMouseDoubleClick( ID string , References []string ) string {
    var Value string = ""
    return Value
}






