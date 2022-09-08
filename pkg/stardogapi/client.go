package stardogapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client holds an HTTPClient and connectivity information
type Client struct {
	BaseURL    string
	Username   string
	Password   string
	HTTPClient *http.Client
}

// errorResponse is an internal struct to decode Stardog error messages
type errorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// Create a new API Client
func NewClient(username, password, baseURL string) *Client {
	return &Client{BaseURL: baseURL,
		Username: username,
		Password: password,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// Send an HTTP request to the Stardog server and write the response to the body
// Decodes the returned error message and includes it in the error
func (c *Client) sendRequest(ctx context.Context, method string, path string, body any, response any) error {
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", c.BaseURL, path), nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Username, c.Password)

	if body != nil {
		pipeReader, pipeWriter := io.Pipe()

		go func() error {
			defer pipeWriter.Close()

			return json.NewEncoder(pipeWriter).Encode(body)
		}()

		req.Body = pipeReader
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error with status code: %d", res.StatusCode)
	}

	if response != nil {
		if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
			return err
		}
	}

	return nil
}
