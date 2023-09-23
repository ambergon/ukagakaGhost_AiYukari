package main

import (
    //"fmt"
    "math/rand"
)



var NadeCount int = 0
func Nade( Char string , Parts string ) string {

    if Char != "0" || Parts == "" {
        return ""
    }

    NadeCount++
    //fmt.Println( NadeCount )
    if NadeCount < 90 {
        return ""
    }
    NadeCount = 0

    res := ""
    if Parts == "Head" {
        i := rand.Intn( len( NadeHeadPrompt ) )
        res = NadeHeadPrompt[ i ]

    } else if Parts == "Bust" {
        i := rand.Intn( len( NadeBustPrompt ) )
        res = NadeBustPrompt[ i ]
    }


    if res != "" {
        go AiTalk( res , "OnMouseMove" )
        return MyChatGPTConfig.NoticeText
    }
    return ""
}


