package main

import "os"

//Errors type
type Errors struct {
	Err []Error `json:"errors"`
}

//Error type
type Error struct {
	Location string `json:"location"`
	Param    string `json:"param"`
	Value    string `json:"value"`
	Msg      string `json:"msg"`
}

const tokenName = "Token"

var tokenSecret = []byte(os.Getenv("JWT_TOKEN"))
