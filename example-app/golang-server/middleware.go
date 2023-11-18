package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

func LogPretty(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		contentType := (strings.Split(request.Header.Get("content-type"), ";"))[0]
		fmt.Println("Content Type of the request:", color.HiYellowString(contentType))
		fmt.Println("HTTP Request Method:", color.HiGreenString(request.Method))
		fmt.Println("URL Being Requested:", color.HiBlueString("http://"+request.Host+request.URL.Path))
		fmt.Println("Protocol In Use:", color.HiBlueString(request.Proto))
		fmt.Println("Request Body Byte Size:", color.HiWhiteString(fmt.Sprint(request.ContentLength)))
		fmt.Println("Request from IP Address:", color.HiMagentaString(request.RemoteAddr))

		// Call the next handler in the chain
		next.ServeHTTP(w, request)
	})
}
