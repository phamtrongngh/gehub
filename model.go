package main

import "net/http"

type ClientResponse struct {
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body"`
}

type Client struct {
	ID      string               `json:"-"`
	Port    uint                 `json:"port"`
	Alias   string               `json:"alias"`
	FwdChan chan *ClientResponse `json:"-"`
}
