package main

import (
    "math/rand"
)

func MouseDoubleClick( Parts string) string {
    res := ""
    if Parts == "Head" {
        i := rand.Intn( len( DoubleClickHeadPrompt ) )
        res = DoubleClickHeadPrompt[ i ]

    } else if Parts == "Bust" {
        i := rand.Intn( len( DoubleClickBustPrompt) )
        res = DoubleClickBustPrompt[ i ]
    }


    if res != "" {
        go AiTalk( res , "OnMouseMove" )
        return MyChatGPTConfig.NoticeText
    }
    return ""
}

