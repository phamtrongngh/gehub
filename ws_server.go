package main

import (
	socketio "github.com/googollee/go-socket.io"
)

func StartWsServer() *socketio.Server {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(c socketio.Conn) error {
		Logger.Infof("Client %s has connected to WebSocket server.", c.ID())
		return nil
	})

	server.OnEvent("/", "expose", func(c socketio.Conn, client *Client) {
		if server.Count() > ConnectionLimit {
			c.Emit("expose", map[string]string{
				"error": "Server is overloaded, please come back later",
			})
			return
		}

		client.ID = c.ID()
		client.FwdChan = make(chan *ClientResponse, 10)

		if client.Alias == "" {
			client.Alias = RandomString(AliasLength)
		}
		if _, exist := Store.ClientByAlias[client.Alias]; exist {
			c.Emit("expose", map[string]string{
				"error": "Alias already exists, please choose a different one",
			})
			return
		}

		c.SetContext(client)
		Store.ClientByAlias[client.Alias] = client

		Logger.Infof(
			"Client %s has exposed port %d with alias %s.",
			c.ID(),
			client.Port,
			client.Alias,
		)

		c.Emit("expose", client)
	})

	server.OnEvent("/", "unexpose", func(c socketio.Conn) {
		client, ok := c.Context().(*Client)
		if ok {
			delete(Store.ClientByAlias, client.Alias)
		}

		Logger.Infof(
			"Client %s has unexposed port %d with alias %s.",
			c.ID(),
			client.Port,
			client.Alias,
		)

		c.Emit("unexpose")
	})

	server.OnEvent("/", "forward", func(c socketio.Conn, res *ClientResponse) {
		Logger.Infof("Received response from client %s with status %d", c.ID(), res.Status)
		client := c.Context().(*Client)
		client.FwdChan <- res
	})

	server.OnDisconnect("/", func(c socketio.Conn, msg string) {
		client, ok := c.Context().(*Client)
		if ok {
			delete(Store.ClientByAlias, client.Alias)
		}
		Logger.Infof("Client %s has disconnected: %s", c.ID(), msg)
	})

	go server.Serve()

	return server
}
