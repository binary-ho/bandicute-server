package connection

import (
	"bandicute-server/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

type SupabaseConnection struct {
	BaseURL string
	ApiKey  string
	Client  *http.Client
}

const BaseEndpoint = "/rest/v1"

func (conn *SupabaseConnection) NewRequest(ctx context.Context, method DML, table Table, query string, body interface{}) (*http.Request, error) {
	if ctx.Err() != nil {
		return nil, fmt.Errorf("context error before creating request: %v", ctx.Err())
	}

	reqBody := strings.NewReader("")
	if body != nil {
		jsonBody, err := json.Marshal(body)

		logger.Info("request json body", logger.Fields{
			"body": string(jsonBody),
		})

		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = strings.NewReader(string(jsonBody))
	}

	url := buildQueryUrl(conn.BaseURL, table, query)
	req, err := http.NewRequestWithContext(ctx, string(method), url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("apikey", conn.ApiKey)
	req.Header.Set("Authorization", "Bearer "+conn.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	logger.Info("new request created", logger.Fields{
		"req header": req.Header,
		"req URL":    req.URL,
	})

	return req, nil
}

func buildQueryUrl(baseUrl string, table Table, query string) string {
	tablePath := strings.Trim(string(table), "/")
	baseUrl = strings.TrimRight(baseUrl, "/") + "/"
	url := baseUrl + strings.TrimPrefix(BaseEndpoint, "/") + "/" + tablePath

	// Add query parameters if they exist
	if query != "" {
		url += "?" + query
	}

	return url
}

func (conn *SupabaseConnection) Do(req *http.Request, v interface{}) error {
	if req.Context().Err() != nil {
		return &fiber.Error{
			Code:    fiber.StatusRequestTimeout,
			Message: "request context cancelled",
		}
	}

	resp, err := conn.Client.Do(req)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("request failed: %v", err),
		}
	}

	defer resp.Body.Close()

	// Handle non-2xx responses
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &fiber.Error{
			Code:    resp.StatusCode,
			Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode),
		}
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return &fiber.Error{
				Code:    fiber.StatusInternalServerError,
				Message: fmt.Sprintf("failed to decode response: %v", err),
			}
		}
	}

	return nil
}
