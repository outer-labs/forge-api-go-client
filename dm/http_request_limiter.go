package dm

import (
	"context"
	"io"
	"net/http"
)

type HttpRequestLimiter interface {
	HttpRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error)
}
