package main

import (
    "fmt"
    "os"
    "encoding/json"
)



func LoadJson(){
    //ChatGPTConfig
	JsonChatGPTConfig, err := os.Open( Directory + "/Config.json")
	if err != nil {
        fmt.Println( err )
	}
	defer JsonChatGPTConfig.Close()
    decoder := json.NewDecoder( JsonChatGPTConfig )
    err     = decoder.Decode( &MyChatGPTConfig )
	if err != nil {
        fmt.Println(  err  )
    }
}






