package http

import (
	"bytes"
	"net/http"
)

func Do(method, username, token, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, token)
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	return client.Do(req)
}

func DoGet(username, token, url string) (*http.Response, error) {
	return Do("GET", username, token, url, []byte{})
}

func DoPost(username, token, url string, body []byte) (*http.Response, error) {
	return Do("POST", username, token, url, body)
}
