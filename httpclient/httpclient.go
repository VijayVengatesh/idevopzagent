package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

// SendGET sends a GET request with optional query parameters
func SendGET(apiURL string, queryParams map[string]string) (*http.Response, error) {
	reqURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}

	// Add query parameters
	q := reqURL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	reqURL.RawQuery = q.Encode()

	resp, err := http.Get(reqURL.String())
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SendPOST sends a POST request with JSON payload to the given endpoint
func SendPOST(apiURL string, payload interface{}) (*http.Response, error) {
	return sendWithBody("POST", apiURL, payload)
}

// SendPUT sends a PUT request with JSON payload to the given endpoint
func SendPUT(apiURL string, payload interface{}) (*http.Response, error) {
	return sendWithBody("PUT", apiURL, payload)
}

// SendDELETE sends a DELETE request with optional JSON payload
func SendDELETE(apiURL string, payload interface{}) (*http.Response, error) {
	if payload == nil {
		req, err := http.NewRequest("DELETE", apiURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		return client.Do(req)
	}
	return sendWithBody("DELETE", apiURL, payload)
}

// sendWithBody is a helper to send POST/PUT/DELETE with body
func sendWithBody(method, apiURL string, payload interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

// ParseJSON parses the response body into the target struct/interface
func ParseJSON(resp *http.Response, target interface{}) error {
	if resp == nil {
		return errors.New("nil response received")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}
