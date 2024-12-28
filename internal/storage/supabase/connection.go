package supabase

import (
	conn "bandicute-server/internal/storage/repository/connection"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Connection struct {
	BaseURL string
	ApiKey  string
	Client  *http.Client
}

const BaseEndpoint = "/rest/v1"

func (r *Connection) NewRequest(ctx context.Context, method conn.DML, table conn.Table, query string, body interface{}) (*http.Request, error) {
	var reqBody *strings.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = strings.NewReader(string(jsonBody))
	}

	req, err := http.NewRequestWithContext(ctx, string(method), buildQueryUrl(r.BaseURL, table, query), reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", r.ApiKey)
	req.Header.Set("Authorization", "Bearer "+r.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	if method == http.MethodPost || method == http.MethodPatch || method == http.MethodDelete {
		req.Header.Set("Prefer", "return=representation")
	}
	return req, nil
}

func buildQueryUrl(baseUrl string, table conn.Table, query string) string {
	queryParam := strings.TrimPrefix(string(table), "/")
	queryParam = strings.TrimSuffix(queryParam, "?")
	if query != "" {
		queryParam += "?" + query
	}
	return baseUrl + BaseEndpoint + queryParam
}

func (r *Connection) Do(req *http.Request, v interface{}) error {
	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}
