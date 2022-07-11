package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"

	"gehub/pkg"
)

func main() {
	wsServer := socketio.NewServer(nil)

	wsServer.OnConnect("/", func(c socketio.Conn) error {
		pkg.Logger.Infof("Client %s has connected to WebSocket server.", c.ID())
		return nil
	})

	wsServer.OnEvent("/", "expose", func(c socketio.Conn, client *pkg.Client) {
		if wsServer.Count() > pkg.ConnectionLimit {
			c.Emit("expose", map[string]string{
				"error": "Server is overloaded, please come back later",
			})
			return
		}

		client.ID = c.ID()
		client.FwdChan = make(chan *pkg.ClientResponse, 10)

		if client.Alias == "" {
			client.Alias = pkg.RandomString(pkg.AliasLength)
		}
		if _, exist := pkg.Store.ClientByAlias[client.Alias]; exist {
			c.Emit("expose", map[string]string{
				"error": "Alias already exists, please choose a different one",
			})
			return
		}

		proxyUrl, err := url.Parse(pkg.ProxyUrl)
		if err != nil {
			c.Emit("expose", map[string]string{
				"error": "Something went wrong, please try again later",
			})
			return
		}
		client.ProxyUrl = proxyUrl.Scheme + "://" + client.Alias + "." + proxyUrl.Host + proxyUrl.Path

		c.SetContext(client)
		pkg.Store.ClientByAlias[client.Alias] = client

		pkg.Logger.Infof(
			"Client %s has exposed port %d with alias %s.",
			c.ID(),
			client.Port,
			client.Alias,
		)

		c.Emit("expose", client)
	})

	wsServer.OnEvent("/", "unexpose", func(c socketio.Conn) {
		client, ok := c.Context().(*pkg.Client)
		if ok {
			delete(pkg.Store.ClientByAlias, client.Alias)
		}

		pkg.Logger.Infof(
			"Client %s has unexposed port %d with alias %s.",
			c.ID(),
			client.Port,
			client.Alias,
		)

		c.Emit("unexpose")
	})

	wsServer.OnEvent("/", "forward", func(c socketio.Conn, res *pkg.ClientResponse) {
		pkg.Logger.Infof("Received response from client %s with status %d", c.ID(), res.Status)
		client := c.Context().(*pkg.Client)

		client.FwdChan <- res
	})

	wsServer.OnDisconnect("/", func(c socketio.Conn, msg string) {
		client, ok := c.Context().(*pkg.Client)
		if ok {
			delete(pkg.Store.ClientByAlias, client.Alias)
		}
		pkg.Logger.Infof("Client %s has disconnected: %s", c.ID(), msg)
	})

	go wsServer.Serve()

	r := gin.Default()
	r.LoadHTMLFiles("web/index.html")
	r.Static("/assets", "web/assets")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/socket.io/*any", gin.WrapH(wsServer))
	r.POST("/socket.io/*any", gin.WrapH(wsServer))

	r.Any("/proxy/*any", func(ctx *gin.Context) {
		alias := strings.Split(ctx.Request.Host, ".")[0]

		client, ok := pkg.Store.ClientByAlias[alias]
		if !ok {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		bodyBytes, _ := ctx.GetRawData()
		body := pkg.BytesToMap(bodyBytes)
		method := ctx.Request.Method
		headers := ctx.Request.Header
		port := client.Port
		path := strings.Join(
			strings.Split(
				ctx.Request.URL.String(), "/",
			)[2:], "/",
		)

		pkg.Logger.Infof(
			"Received request with alias=%s, method=%s, path=%s",
			alias, method, path,
		)

		wsServer.BroadcastToRoom("/", client.ID, "forward", map[string]any{
			"method":  method,
			"headers": headers,
			"body":    body,
			"path":    path,
			"port":    port,
		})

		res := <-client.FwdChan

		for k, v := range res.Headers {
			ctx.Writer.Header().Set(k, strings.Join(v, ","))
		}

		ctx.Writer.WriteString(res.Body)
		ctx.Writer.WriteHeader(res.Status)

		pkg.Logger.Infof(
			"Forward response with status=%d, headers=%v",
			res.Status, res.Headers,
		)
	})

	if err := r.Run(
		strings.Split(pkg.WsUrl, "://")[1],
	); err != nil {
		panic(err)
	}
}
