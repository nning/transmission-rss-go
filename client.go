package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type RequestArguments map[string]interface{}

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

	url := config.ServerURL()

	request, err := http.NewRequest("GET", url, nil)
	panicOnError(err)

	request.SetBasicAuth(config.Login.Username, config.Login.Password)
	response, err := client.Do(request)
	panicOnError(err)

	_, err = ioutil.ReadAll(response.Body)
	panicOnError(err)

	if response.StatusCode != 409 {
		status := strconv.Itoa(response.StatusCode)
		fmt.Println("SESSION_ID ERROR", status)
		panic("Could not obtain session ID, got HTTP response code " + status + ".")
	}

	sessionId := response.Header["X-Transmission-Session-Id"][0]

	fmt.Println("SESSION_ID", sessionId)

	return sessionId
}

func (self *Client) UpdateSessionId() {
	self.SessionId = getSessionId(self.Config)
}

func (self *Client) rpc(requestBody RequestBody) http.Response {
	client := &http.Client{}

	url := self.Config.ServerURL()

	jsonData, err := json.Marshal(requestBody)
	panicOnError(err)

	request, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	panicOnError(err)

	login := self.Config.Login
	request.SetBasicAuth(login.Username, login.Password)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Transmission-Session-Id", self.SessionId)

	response, err := client.Do(request)
	panicOnError(err)

	// TODO Catch 409, update SessionId, retry
	//      https://github.com/transmission/transmission/blob/master/extras/rpc-spec.txt#L56

	return *response
}

func (self *Client) AddTorrent(link string) (id int, err error) {
	var requestBody RequestBody
	requestBody.Method = "torrent-add"
	requestBody.Arguments = make(map[string]interface{})
	requestBody.Arguments["filename"] = link

	if self.Config.Paused {
		requestBody.Arguments["paused"] = true
	}

	response := self.rpc(requestBody)
	if response.StatusCode != 200 {
		return 0, fmt.Errorf("RPC call error status: %d", response.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	jsonBody := buf.Bytes()

	var jsonResult ResponseBody
	json.Unmarshal(jsonBody, &jsonResult)
	if jsonResult.Result != "ok" {
		return 0, errors.New(jsonResult.Result)
	}

	id = jsonResult.Arguments.TorrentAdded.Id
	if id == 0 {
		id = jsonResult.Arguments.TorrentDuplicate.Id
	}

	return id, nil
}

func (self *Client) SetTorrent(arguments RequestArguments) {
	var requestBody RequestBody
	requestBody.Method = "torrent-set"
	requestBody.Arguments = arguments

	self.rpc(requestBody)
}
