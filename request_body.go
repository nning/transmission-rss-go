package main

type RequestBody struct {
	Method    string `json:"method"`
	Arguments struct {
		Filename string `json:"filename"`
	} `json:"arguments"`
}
