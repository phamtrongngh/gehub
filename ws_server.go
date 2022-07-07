package main

import (
	"errors"

	socketio "github.com/googollee/go-socket.io"
)

func disconnect(c socketio.Conn, err error) {
	client := c.Context().(*Client)
	delete(Store.ClientByAlias, client.Alias)
	c.Emit("disconnect", err.Error())
}

func StartWsServer() *socketio.Server {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(c socketio.Conn) error {
		if server.Count() > ConnectionLimit {
			err := errors.New("concurrent connections limit exceeded")
			disconnect(c, err)
			return err
		}

		url := c.URL()

		client := &Client{
			ID:    c.ID(),
			Port:  url.Query().Get("port"),
			Alias: url.Query().Get("alias"),
		}

		if client.Alias == "" {
			client.Alias = RandomString(AliasLength)
		}
		if _, exist := Store.ClientByAlias[client.Alias]; exist {
			err := errors.New("alias already exists")
			disconnect(c, err)
			return err
		}

		c.SetContext(client)
		Store.ClientByAlias[client.Alias] = client

		Logger.Infof(
			"Client '%s' has connected and exposed port '%s' with alias '%s'.",
			c.ID(),
			client.Port,
			client.Alias,
		)

		return nil
	})

	server.OnEvent("/", "info", func(c socketio.Conn) {
		client := c.Context().(*Client)
		c.Emit("info", client)
	})

	server.OnEvent("/", "forward", func(c socketio.Conn, res *ClientResponse) {
		Logger.Infof("Received response from client '%s' with status %d", c.ID(), res.Status)

	})

	server.OnDisconnect("/", func(c socketio.Conn, s string) {
		disconnect(c, nil)
		Logger.Infof("Client '%s' has disconnected.", c.ID())
	})

	go server.Serve()

	return server
}
