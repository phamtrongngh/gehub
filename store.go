package main

type store struct {
	ClientByAlias map[string]*Client
}

func NewStore() *store {
	return &store{
		ClientByAlias: make(map[string]*Client),
	}
}
