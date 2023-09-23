package main

import (
    "io/ioutil"
    "os"
    "bufio"
)


var DefaultPrompt           string
var WheelClickMousePrompt   []string
var NadeHeadPrompt          []string
var NadeBustPrompt          []string
var DoubleClickHeadPrompt   []string
var DoubleClickBustPrompt   []string


func LoadPrompt(){
    DefaultPrompt = ""

    bytes, err := ioutil.ReadFile( Directory + "/DefaultPrompt.txt" )
    if err != nil { panic(err) }
    DefaultPrompt = string(bytes)


    fp, err := os.Open( Directory + "/Nade/HeadPrompt.txt" )
    if err != nil { panic( err ) }
    defer fp.Close()
    scanner := bufio.NewScanner(fp)
    for scanner.Scan() {
        if scanner.Text() == "" { continue }
        NadeHeadPrompt = append( NadeHeadPrompt , scanner.Text())
    }
    if err = scanner.Err(); err != nil { panic( err ) }


    fp, err = os.Open( Directory + "/Nade/BustPrompt.txt" )
    if err != nil { panic( err ) }
    defer fp.Close()
    scanner = bufio.NewScanner(fp)
    for scanner.Scan() {
        if scanner.Text() == "" { continue }
        NadeBustPrompt = append( NadeBustPrompt , scanner.Text())
    }
    if err = scanner.Err(); err != nil { panic( err ) }


    fp, err = os.Open( Directory + "/WheelClick/MousePrompt.txt" )
    if err != nil { panic( err ) }
    defer fp.Close()
    scanner = bufio.NewScanner(fp)
    for scanner.Scan() {
        if scanner.Text() == "" { continue }
        WheelClickMousePrompt = append( WheelClickMousePrompt , scanner.Text())
    }
    if err = scanner.Err(); err != nil { panic( err ) }


    fp, err = os.Open( Directory + "/DoubleClick/HeadPrompt.txt" )
    if err != nil { panic( err ) }
    defer fp.Close()
    scanner = bufio.NewScanner(fp)
    for scanner.Scan() {
        if scanner.Text() == "" { continue }
        DoubleClickHeadPrompt = append( DoubleClickHeadPrompt , scanner.Text())
    }
    if err = scanner.Err(); err != nil { panic( err ) }


    fp, err = os.Open( Directory + "/DoubleClick/BustPrompt.txt" )
    if err != nil { panic( err ) }
    defer fp.Close()
    scanner = bufio.NewScanner(fp)
    for scanner.Scan() {
        if scanner.Text() == "" { continue }
        DoubleClickBustPrompt = append( DoubleClickBustPrompt , scanner.Text())
    }
    if err = scanner.Err(); err != nil { panic( err ) }

}
