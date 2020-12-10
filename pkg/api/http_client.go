package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type httpClient struct {
	client  *http.Client
	baseURL string

	baseHeaders http.Header
}

func NewHTTPClient(baseURL string) *httpClient {
	baseHeaders := make(http.Header)

	baseHeaders.Set("Accept", "application/json")

	return &httpClient{
		client: &http.Client{
			Timeout: requestTimeout,
		},
		baseURL:     baseURL,
		baseHeaders: baseHeaders,
	}
}

func (c *httpClient) SetAuthToken(token string) {
	c.baseHeaders.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func (c *httpClient) GET(ctx context.Context, path string) ([]byte, error) {
	return c.do(ctx, http.MethodGet, path, nil)
}

func (c *httpClient) POST(ctx context.Context, path string, body interface{}) ([]byte, error) {
	bodyBuf := new(bytes.Buffer)

	if err := json.NewEncoder(bodyBuf).Encode(body); err != nil {
		return nil, fmt.Errorf("body encode to json failed: %w", err)
	}

	return c.do(ctx, http.MethodPost, path, bodyBuf)
}

func (c *httpClient) checkInternetConnection() error {
	res, err := c.client.Get(c.baseURL + return204Path)

	if err != nil {
		return &NoInternetError{
			Err: err,
		}
	}

	_ = res.Body.Close()

	return nil
}

func (c *httpClient) do(ctx context.Context, method, path string, body io.Reader) ([]byte, error) {
	if err := c.checkInternetConnection(); err != nil {
		return nil, err
	}

	u := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, u, body)

	if err != nil {
		return nil, fmt.Errorf("http request create failed: %w", err)
	}

	// set base headers
	req.Header = c.baseHeaders.Clone()

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("http do error: %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if err := c.checkResponse(res); err != nil {
		return nil, err
	}

	var response struct {
		Data json.RawMessage `json:"data"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("res body unmarshall fail: %w", err)
	}

	return response.Data, nil
}

func (c *httpClient) checkResponse(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}

	if res.Header.Get("Content-Type") == "application/json" {
		switch res.StatusCode {
		default:
			var resErr struct {
				Error RequestError `json:"error"`
			}

			if err := json.NewDecoder(res.Body).Decode(&resErr); err != nil {
				return fmt.Errorf("error body unmarshall fail: %w", err)
			}

			return &resErr.Error

		case http.StatusUnauthorized:
			var resErr struct {
				Error UnauthorizedError `json:"error"`
			}

			if err := json.NewDecoder(res.Body).Decode(&resErr); err != nil {
				return fmt.Errorf("error body unmarshall fail: %w", err)
			}

			return &resErr.Error

		case http.StatusBadRequest:
			var resErr struct {
				Error ValidationError `json:"error"`
			}

			if err := json.NewDecoder(res.Body).Decode(&resErr); err != nil {
				return fmt.Errorf("error body unmarshall fail: %w", err)
			}

			return &resErr.Error
		}
	}

	// if response isn't json
	bodyRaw, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return &RequestError{
			StatusCode: res.StatusCode,
		}
	}

	return &RequestError{
		StatusCode: res.StatusCode,
		Message:    string(bodyRaw),
	}
}
