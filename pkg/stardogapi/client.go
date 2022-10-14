package stardogapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
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

// Send an HTTP request to the Stardog server and decode the JSON response (incl. JSON errors)
func (c *Client) sendRequest(ctx context.Context, method string, path string, body any, response any) error {
	bodyBuffer := &bytes.Buffer{}

	err := json.NewEncoder(bodyBuffer).Encode(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", c.BaseURL, path), bodyBuffer)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)

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

// Send a multipart HTTP request to the Stardog server and decode the JSON response (incl. JSON errors)
func (c *Client) sendMultipartJsonRequest(ctx context.Context, method string, path string, body map[string]any, response any) error {
	bodyBuffer := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(bodyBuffer)

	for k, v := range body {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, k))
		h.Set("Content-Type", "application/json")
		field, err := multipartWriter.CreatePart(h)
		if err != nil {
			return err
		}

		err = json.NewEncoder(field).Encode(v)
		if err != nil {
			return err
		}
	}

	err := multipartWriter.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", c.BaseURL, path), bodyBuffer)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return fmt.Errorf("%s: %s", errRes.Code, errRes.Message)
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
