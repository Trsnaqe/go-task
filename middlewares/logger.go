package middlewares

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestAsText, err := getRequestAsText(r)
		if err != nil {
			log.Printf("Error logging request: %v", err)
		} else {
			log.Print(requestAsText)
		}
		next.ServeHTTP(w, r)
	})
}

func getRequestAsText(r *http.Request) (string, error) {
	rawRequestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("error reading request body: %v", err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(rawRequestBody))

	requestText := fmt.Sprintf("\n%s------ Request Details ------%s\n", ColorBlue, ColorReset)
	requestText += fmt.Sprintf("%sMethod:%s %s\n", ColorYellow, ColorReset, r.Method)
	requestText += fmt.Sprintf("%sURL:%s %s://%s%s%s\n", ColorYellow, ColorReset, r.URL.Scheme, r.Host, r.URL.Path, r.URL.RawQuery)
	/*
		requestText += fmt.Sprintf("\n%s------ Headers ------%s\n", ColorYellow, ColorReset)

		for key, values := range r.Header {
			requestText += fmt.Sprintf("%s%s:%s %s\n", ColorGreen, key, ColorReset, values)

		}
	*/
	requestText += fmt.Sprintf("\n%s------ Content ------%s\n", ColorYellow, ColorReset)

	requestText += fmt.Sprintf("\n%sContent:%s\n%s\n", ColorYellow, ColorReset, rawRequestBody)

	return requestText, nil
}
