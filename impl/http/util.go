package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Request is used to ensure that no leaky code (unclosed memory/goroutines)
// and that there is a sole function for handling the actual HTTP requests
func Request(ctx context.Context, methodType string, res interface{}, body interface{}, url url.URL) (*http.Response, error) {
	var (
		bodyBytes []byte
		err       error
	)

	if body != nil {
		bodyBytes, err = json.Marshal(&body)
		if err != nil {
			// TODO: wrapping and logging
			// log
			return nil, err
		}
	}

	req, err := http.NewRequest(methodType, url.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		// TODO: wrapping and logging
		// log
		return nil, err
	}

	applyHeaders(req)

	httpRes, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		// TODO: wrapping and logging
		return nil, err
	}

	defer closeResBody(httpRes)

	httpBody, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		// TODO: wrapping and logging
		return httpRes, err
	}

	if httpBody == nil {
		// TODO: logging
		return httpRes, nil
	}

	err = json.Unmarshal(httpBody, res)
	if err != nil {
		// TODO: wrapping and logging
		return httpRes, err
	}

	return httpRes, nil
}

func closeResBody(res *http.Response) {
	var err = res.Body.Close()
	if err != nil {
		// TODO: wrapping and logging
	}
}

func applyHeaders(req *http.Request) {
	// Apply headers inside of here
}
