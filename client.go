package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RequestArguments struct {
	Filename       string  `json:"filename"`
	Paused         bool    `json:"paused"`
	Ids            []int   `json:"ids"`
	SeedRatioLimit float32 `json:"seedRatioLimit"`
	SeedRatioMode  int     `json:"seedRatioMode"`
}

type RequestBody struct {
	Method    string           `json:"method"`
	Arguments RequestArguments `json:"arguments"`
}

type ResponseBody struct {
	Result    string `json:"result"`
	Arguments struct {
		TorrentAdded struct {
			Id int `json:"id"`
		} `json:"torrent-added"`
		TorrentDuplicate struct {
			Id int `json:"id"`
		} `json:"torrent-duplicate"`
	} `json:"arguments"`
}

type Client struct {
	Config    *Config
	SessionId string
}

func NewClient(config *Config) *Client {
	client := Client{
		Config:    config,
		SessionId: getSessionId(config),
	}

	return &client
}

func getSessionId(config *Config) string {
	client := &http.Client{}

	url := getUrl(config)

	request, err := http.NewRequest("GET", url, nil)
	panicOnError(err)

	request.SetBasicAuth(config.Login.Username, config.Login.Password)
	response, err := client.Do(request)
	panicOnError(err)

	_, err = ioutil.ReadAll(response.Body)
	panicOnError(err)

	sessionId := response.Header["X-Transmission-Session-Id"][0]

	fmt.Println("SESSION_ID " + sessionId)

	return sessionId
}

func (self *Client) UpdateSessionId() {
	self.SessionId = getSessionId(self.Config)
}

func (self *Client) rpc(requestBody RequestBody) http.Response {
	client := &http.Client{}

	url := getUrl(self.Config)

	jsonData, err := json.Marshal(requestBody)
	panicOnError(err)

	// fmt.Println(string(jsonData))

	request, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	panicOnError(err)

	login := self.Config.Login
	request.SetBasicAuth(login.Username, login.Password)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Transmission-Session-Id", self.SessionId)

	response, err := client.Do(request)
	panicOnError(err)

	return *response
}

func (self *Client) AddTorrent(link string) int {
	var requestBody RequestBody
	requestBody.Method = "torrent-add"
	requestBody.Arguments.Filename = link

	if self.Config.Paused {
		requestBody.Arguments.Paused = true
	}

	response := self.rpc(requestBody)

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	jsonBody := buf.Bytes()

	var jsonResult ResponseBody
	json.Unmarshal(jsonBody, &jsonResult)

	id := jsonResult.Arguments.TorrentAdded.Id
	if id == 0 {
		id = jsonResult.Arguments.TorrentDuplicate.Id
	}

	return id
}

func (self *Client) SetTorrent(arguments RequestArguments) {
	var requestBody RequestBody
	requestBody.Method = "torrent-set"
	requestBody.Arguments = arguments

	// fmt.Println(requestBody)

	self.rpc(requestBody)

	// response := self.rpc(requestBody)
	//
	// buf := new(bytes.Buffer)
	// buf.ReadFrom(response.Body)
	//
	// fmt.Println(buf.String())
}
