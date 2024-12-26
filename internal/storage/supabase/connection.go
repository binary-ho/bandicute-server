package supabase

import (
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

func (r *Connection) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	var reqBody *strings.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = strings.NewReader(string(jsonBody))
	}

	req, err := http.NewRequestWithContext(ctx, method, r.BaseURL+BaseEndpoint+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", r.ApiKey)
	req.Header.Set("Authorization", "Bearer "+r.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
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
