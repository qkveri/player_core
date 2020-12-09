package api

import "context"

type Client interface {
	SetAuthToken(token string)

	GET(ctx context.Context, path string) ([]byte, error)
	POST(ctx context.Context, path string, body interface{}) ([]byte, error)
}
