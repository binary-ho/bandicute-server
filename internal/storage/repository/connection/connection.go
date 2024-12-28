package connection

import (
	"bandicute-server/internal/storage/supabase"
	"context"
	"net/http"
	"strings"
)

type Connection interface {
	NewRequest(ctx context.Context, method DML, table Table, query string, body interface{}) (*http.Request, error)
	Do(req *http.Request, v interface{}) error
}

func NewConnection(baseURL, apiKey string) Connection {
	return &supabase.Connection{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		ApiKey:  apiKey,
		Client:  &http.Client{},
	}
}

type DML string
type Query string
type Table string
