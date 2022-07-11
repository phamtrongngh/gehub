package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"gehub/pkg"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	wsServerUrl, err := url.Parse(pkg.WsUrl)
	if err != nil {
		panic(err)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(wsServerUrl)
	reverseProxy.Director = func(req *http.Request) {
		req.URL.Scheme = wsServerUrl.Scheme
		req.URL.Host = wsServerUrl.Host
		req.URL.Path = "/proxy" + req.URL.Path
	}

	reverseProxy.ModifyResponse = func(resp *http.Response) error {
		// check if response is a image
		if strings.HasPrefix(resp.Header.Get("Content-Type"), "image") {
			resp.Header.Set("Cache-Control", "public, max-age=31536000")
		}

		return nil
	}

	server.Any("/*any", gin.WrapH(reverseProxy))

	if err := server.Run(
		strings.Split(pkg.ProxyUrl, "://")[1],
	); err != nil {
		panic(err)
	}
}
