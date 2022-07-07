package main

type Client struct {
	ID    string `json:"id"`
	Port  string `json:"port"`
	Alias string `json:"alias"`
}

type ClientResponse struct {
	Status  int            `json:"status"`
	Headers map[string]any `json:"headers"`
	Body    any            `json:"body"`
}
