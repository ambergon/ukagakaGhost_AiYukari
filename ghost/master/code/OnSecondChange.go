
package main


import (
    "fmt"
    "os"
)

var NextSecondTalk []string 
var OnSaveTalkScript string
//APIの補充カウント
var ChargeAPISec    int = 0

func OnSecondChange( NOTIFY bool ) string {
    var Value string = ""
    if NOTIFY == false {
        if len( NextSecondTalk ) != 0 {
            Value  = NextSecondTalk[0]
            //AIのトークなので保存できるように置いておく。
            OnSaveTalkScript = NextSecondTalk[0]
            NextSecondTalk   = NextSecondTalk[1:]
        }
    }


    //ChargeAPISecMaxごとに補充する。
    if MyChatGPTConfig.ChargeAPIMax > MyChatGPTConfig.ChargeAPI {
        ChargeAPISec++
        if ChargeAPISec >= MyChatGPTConfig.ChargeAPISecMax {
            MyChatGPTConfig.ChargeAPI++
            ChargeAPISec = 0
        }
    }
    return Value
}


func SaveTalk( SaveText string ){
    fp,err := os.OpenFile( Directory + "/log/log.txt" , os.O_WRONLY|os.O_APPEND|os.O_CREATE , 0665 )
    if err != nil {
        fmt.Println( "保存error" )
        return 
    }
    defer fp.Close()
    fmt.Fprintln( fp , SaveText )
}























