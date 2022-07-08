package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func SetupHttpServer(wsServer *socketio.Server) *gin.Engine {
	r := gin.Default()

	r.LoadHTMLFiles("public/index.html")
	r.Static("/assets", "public/assets")
	r.GET("/socket.io/*any", gin.WrapH(wsServer))
	r.POST("/socket.io/*any", gin.WrapH(wsServer))

	r.GET("/connect", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.Any("/*any", func(ctx *gin.Context) {
		alias := strings.Split(ctx.Request.Host, ".")[0]

		client, ok := Store.ClientByAlias[alias]
		if !ok {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		bodyBytes, _ := ctx.GetRawData()
		body := BytesToMap(bodyBytes)
		method := ctx.Request.Method
		headers := ctx.Request.Header
		port := client.Port
		path := strings.Join(
			strings.Split(
				ctx.Request.URL.String(), "/",
			)[1:], "/",
		)

		Logger.Infof(
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

		// // add client response headers
		for k, v := range res.Headers {
			fmt.Println("key", k)
			ctx.Writer.Header().Set(k, strings.Join(v, ","))
		}

		ctx.Writer.WriteString(res.Body)
		ctx.Writer.WriteHeader(res.Status)

		Logger.Infof(
			"Forward response with status=%d, headers=%v",
			res.Status, res.Headers,
		)
	})

	if err := r.Run(fmt.Sprintf(":%s", Port)); err != nil {
		Logger.Fatalln(err)
	}

	return r
}
