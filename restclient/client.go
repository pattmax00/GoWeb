package restclient

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
)

// SendRequest sends an HTTP request to a URL and includes the specified headers and body.
// A body can be nil for GET requests, a map[string]string for multipart/form-data requests,
// or a struct for JSON requests
func SendRequest(url string, method string, headers map[string]string, body interface{}) (http.Response, error) {
	var reqBody *bytes.Buffer
	var contentType string

	switch v := body.(type) {
	case nil:
		reqBody = bytes.NewBuffer([]byte(""))
	case map[string]string:
		reqBody = &bytes.Buffer{}
		writer := multipart.NewWriter(reqBody)
		for key, value := range v {
			err := writer.WriteField(key, value)
			if err != nil {
				return http.Response{}, err
			}
		}

		err := writer.Close()
		if err != nil {
			return http.Response{}, err
		}

		contentType = writer.FormDataContentType()
	default:
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return http.Response{}, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
		contentType = "application/json"
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return http.Response{}, err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http.Response{}, err
	}

	return *resp, nil
}
