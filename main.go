package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var Store = NewStore()

func main() {
	r := gin.Default()

	wsServer := StartWsServer()
	defer wsServer.Close()

	r.LoadHTMLGlob("public/*.html")
	r.Static("/assets", "public/assets")

	r.GET("/socket.io/*any", gin.WrapH(wsServer))
	r.POST("/socket.io/*any", gin.WrapH(wsServer))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/test", func(c *gin.Context) {
		go func() {
			time.Sleep(2 * time.Second)
			c.JSON(200, "hello")
		}()
	})

	r.Any("/:alias/*any", func(c *gin.Context) {
		alias := c.Param("alias")

		if alias == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		client, ok := Store.ClientByAlias[alias]
		if !ok {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		bodyBytes, _ := c.GetRawData()
		body := BytesToMap(bodyBytes)
		method := c.Request.Method
		headers := c.Request.Header
		path := strings.Join(
			strings.Split(
				c.Request.URL.String(), "/",
			)[2:], "/",
		)
		port := client.Port

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

		c.JSON(http.StatusOK, nil)
	})

	if err := r.Run(fmt.Sprintf(":%s", Port)); err != nil {
		Logger.Fatalln(err)
	}
}
