package main

/*
   #include <windows.h>
   #include <stdlib.h>
   #include <string.h>
*/
import "C"

import (
    "fmt"
    "unsafe"
    "strings"
    "regexp"
    //"strconv"
)


func main() {
    fmt.Println( "test" )
}


var Directory string

var References []string 
var CheckID         = regexp.MustCompile("^ID: ")
var CheckReference  = regexp.MustCompile("^Reference.+?: ")

var AutoTalkSwitch bool
var WordList []string

type ResponseStruct struct {
    Shiori  string
    Sender  string
    Charset string
    Marker  string
    Value   string
}
func GetResponse( r *ResponseStruct ) string {
    V := ""
    M := ""
    if r.Value  != "" { V = "Value: "  + r.Value     + "\r\n" }
    if r.Marker != "" { M = "Marker: " + r.Marker    + "\r\n" }
    res :=  r.Shiori    + "\r\n" + 
            r.Sender    + "\r\n" + 
            r.Charset   + "\r\n" + 
            M + V + "\r\n\r\n"
    return res
}


//export load
func load(h C.HGLOBAL, length C.long ) C.BOOL {
    fmt.Println( "call load golang shiori" )
    Directory = C.GoStringN(( *C.char )( unsafe.Pointer( h )), ( C.int )( length ))
    fmt.Println( Directory  )

    //設定読み込み。
    LoadJson()

    //プロンプト読み込み
    LoadPrompt()

    AutoTalkSwitch = false

	C.GlobalFree( h )
	return C.TRUE
}


//export unload
func unload() bool {
    fmt.Println( "call unload golang shiori" )
	return true
}


//export request
func request( h C.HGLOBAL, length *C.long ) C.HGLOBAL {
	RequestText := C.GoStringN(( *C.char )( unsafe.Pointer( h )), ( C.int )( *length ))
	//RequestText := C.GoStringN(( *C.char )( unsafe.Pointer( h )), (( C.int )( *length ) - 1))
	C.GlobalFree( h )


    Value           := ""
    Marker          := ""
    ID              := ""
    References      = []string{}
    var NOTIFY bool = false

    //am  := time.Now().Format( "pm" )
    //MonthStr  := time.Now().Format( "01" )
    //DayStr    := time.Now().Format( "02" )
    //HourStr  := time.Now().Format( "15" )
    //MinStr  := time.Now().Format( "04" )
    //Hour ,_  := strconv.Atoi( HourStr ) 
    //Min ,_  := strconv.Atoi( MinStr ) 

    Response := new( ResponseStruct )
    Response.Sender  = "Sender: GolangAI"
    Response.Charset = "Charset: UTF-8"

    //IDとReference
    //必要な情報を分解する。
    RequestLines := strings.Split( RequestText , "\r\n" )
    for _ , line := range RequestLines {
        if( line == "NOTIFY SHIORI/3.0" ){
            NOTIFY = true

        } else if CheckID.MatchString( line )  {
            //fmt.Println( line )
            ID = CheckID.ReplaceAllString( line , "" )

        } else if CheckReference.MatchString( line )  {
            //fmt.Println( line )
            ref := CheckReference.ReplaceAllString( line , "" )
            References = append( References , ref )

        } else {
            //fmt.Println( line )
        }
    }




    //実行関数
    if ID == "OnSecondChange" {
        Value = OnSecondChange( NOTIFY )

    } else if ID == "OnMinuteChange"  {
        if AutoTalkSwitch == true{
            go TopicTalk( ID )
        }

    } else if ID == "OnHourTimeSignal" {


    } else if ID == "OnSurfaceRestore"  {

    } else if ID == "OnMouseMove" {
        Value = Nade( References[3] , References[4] )

    } else if ID == "OnMouseWheel" {
        Value = OnMouseWheel( References[ 2 ] )

    } else if ID == "OnMouseClick" {
        if References[3] == "0" && References[5] == "2" {
            Value = WheelClick( References[4] )
        }

    } else if ID == "OnMouseDoubleClick" {
        if References[ 4 ] == "" {
            AiDelete()
            Value = "削除しました。"
        } else {
            Value = MouseDoubleClick( References[ 4 ] )
        }

        //Value = OnMouseDoubleClick( ID , References )
    } else if ID == "OnKeyPress" {


    ////fileドロップ回り
    ////複数ファイルがbyte1で区切られてくる。
    //} else if ID == "OnFileDropEx"  {
    //    FilesLen    := 0
    //    FileText    := References[0]
    //    Files       := strings.Split( FileText , "\u0001" )
    //    for( FilesLen != len( Files )) {
    //        var ReadFile = regexp.MustCompile(".txt$")
    //        if ReadFile.MatchString( Files[ FilesLen ] )  {
    //            fmt.Println( "処理ファイル : " + Files[ FilesLen ] )
    //        } else {
    //            fmt.Println( "スキップ : " + Files[ FilesLen ] )
    //        }
    //        FilesLen++
    //    }


    //対話
    } else if ID == "OnCommunicate"  {
        //userからのコメントならば。
        if References[ 0 ] == "User" || References[ 0 ] == "user" {
            fmt.Print( "\n>>  " )
            fmt.Println( References[ 1 ]  )
            go AiTalk( References[ 1 ] , ID )
            Value = MyChatGPTConfig.NoticeText
        }
    } else if ID == "OnCommunicateBoxOpen" {
        Value  = "\\0\\![open,communicatebox]" 

    } else if ID == "OnBoot" {
        Value = "\\1\\s[-1]\\0\\![open,communicatebox]"

    } else if ID == "OnClose"  {
        Value = "\\w9\\w9\\w9\\w9\\w9\\w9\\-\\e"

    } else if ID == "OnTranslate"  {
        Value = strings.Replace( References[0] , "。" , "。\\n\\w9" , -1 )
        //Google検索できるようにしてもいいが、パーセントエンコーディングする必要があるな。


    } else if ID == "OnBalloonTimeout"  {
    } else if ID == "rateofusegraph"  {
    } else if ID == "OnBatteryNotify"  {
    } else if ID == "OnOSUpdateInfo"  {
    } else {
        fmt.Println( "no touch :" + ID )
        //n := 0
        //for( len( References ) != n ){
        //    fmt.Println( "no touch :" + References[n] )
        //    n++
        //}
    }


    if Value == "" {
        Response.Shiori  = "SHIORI/3.0 204 No Content"
    } else {
        Response.Shiori = "SHIORI/3.0 200 OK"
        Response.Value  = Value
    }
    if Marker != "" {
        Response.Marker  = Marker
    }

    res_buf := C.CString( GetResponse( Response ))
    defer C.free( unsafe.Pointer( res_buf ))

	res_size := C.strlen( res_buf )
	ret      := C.GlobalAlloc( C.GPTR , ( C.SIZE_T )( res_size ))
	C.memcpy(( unsafe.Pointer )( ret ) , ( unsafe.Pointer )( res_buf ) , res_size )
	*length = ( C.long )( res_size )
	return ret
}



















