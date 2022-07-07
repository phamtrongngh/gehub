package main

import (
	"encoding/json"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func BytesToMap(bytes []byte) (body map[string]any) {
	json.Unmarshal(bytes, &body)
	return
}

// var ForbiddenHeaders = []string{
// 	"X-Requested-With",
// 	"X-Requested-without",
// 	"User-Agent",
// 	"Accept-Charset",
// 	"Accept-Encoding",
// 	"Access-Control-Request-Headers",
// 	"Access-Control-Request-Method",
// 	"Connection",
// 	"Content-Length",
// 	"Cookie",
// 	"Date",
// 	"DNT",
// 	"Expect",
// 	"Feature-Policy",
// 	"Host",
// 	"Keep-Alive",
// 	"Origin",
// 	"Proxy-",
// 	"Sec-",
// 	"Referer",
// 	"TE",
// 	"Trailer",
// 	"Transfer-Encoding",
// 	"Upgrade",
// 	"Via",
// }

// func FilterHeaders(headers http.Header) (filterHeaders http.Header) {
// 	filterHeaders = make(http.Header)
// 	for key, val := range headers {
// 		isValid := true
// 		for _, fh := range ForbiddenHeaders {
// 			if strings.Contains(key, fh) {
// 				isValid = false
// 			}
// 		}
// 		if isValid {
// 			filterHeaders[key] = val
// 		}
// 	}
// 	return
// }
