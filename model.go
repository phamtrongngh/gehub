package main

type Client struct {
	ID    string `json:"id"`
	Port  string `json:"port"`
	Alias string `json:"alias"`
}

type ForwardRequest struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Header any    `json:"header"`
	Body   any    `json:"body"`
	Port   string `json:"port"`
}
