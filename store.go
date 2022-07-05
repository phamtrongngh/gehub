package main

type appStore struct {
	ClientByAlias map[string]*Client
}

var Store *appStore

func init() {
	Store = &appStore{
		ClientByAlias: make(map[string]*Client),
	}
}
