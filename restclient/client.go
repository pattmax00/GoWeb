package restclient

import "net/http"

// SendRequest sends an HTTP request to a URL and includes the specified headers and body
func SendRequest(url string, method string, headers map[string]string, body map[string]string) (http.Response, error) {
	// Create request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return http.Response{}, err
	}

	// Add headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Add body
	for key, value := range body {
		req.Form.Add(key, value)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http.Response{}, err
	}

	return *resp, nil
}
