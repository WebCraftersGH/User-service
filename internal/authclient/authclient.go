package authclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrUnauthorized = errors.New("unauthorized")

type CheckResponse struct {
	UserID string `json:"user_id"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Check(ctx context.Context, token string) (uuid.UUID, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return uuid.Nil, ErrUnauthorized
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+"/auth/check",
		nil,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("authclient: build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return uuid.Nil, fmt.Errorf("authclient: do request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var out CheckResponse
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return uuid.Nil, fmt.Errorf("authclient: decode response: %w", err)
		}

		userID, err := uuid.Parse(out.UserID)
		if err != nil {
			return uuid.Nil, fmt.Errorf("authclient: parse user_id: %w", err)
		}

		return userID, nil

	case http.StatusUnauthorized, http.StatusForbidden:
		return uuid.Nil, ErrUnauthorized

	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return uuid.Nil, fmt.Errorf(
			"authclient: unexpected status=%d body=%s",
			resp.StatusCode,
			strings.TrimSpace(string(body)),
		)
	}
}
