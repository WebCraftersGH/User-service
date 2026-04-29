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

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/WebCraftersGH/User-service/pkg/logging"
)

var ErrUnauthorized = errors.New("unauthorized")

type AuthResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     logging.Logger
}

func New(baseURL string, logger logging.Logger) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		logger: logger,
	}
}

func (c *Client) Check(ctx context.Context, token string) (uuid.UUID, error) {
	token = strings.TrimSpace(token)
	c.logger.WithField("token", token).Info("check auth start")
	if token == "" {
		return uuid.Nil, ErrUnauthorized
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
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

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return uuid.Nil, fmt.Errorf("authclient: decode response: %w", err)
	}

	c.logger.WithField("check_response", authResp).Info("get auth struct")

	switch resp.StatusCode {
	case http.StatusOK:
		// Если успешный ответ, просто парсим JWT и извлекаем uid (без верификации)
		if authResp.Error != "" {
			return uuid.Nil, ErrUnauthorized
		}
		
		userID, err := c.extractUserIDFromToken(token)
		if err != nil {
			return uuid.Nil, fmt.Errorf("authclient: extract user_id from token: %w", err)
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

func (c *Client) extractUserIDFromToken(tokenString string) (uuid.UUID, error) {
	// Парсим JWT без верификации (токен уже проверен auth-service)
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse token: %w", err)
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid claims type")
	}
	
	// Извлекаем uid из claims (как показано в логе)
	var uidStr string
	
	// Пробуем разные варианты названий полей
	if uid, exists := claims["uid"]; exists {
		uidStr = fmt.Sprintf("%v", uid)
	} else if userID, exists := claims["user_id"]; exists {
		uidStr = fmt.Sprintf("%v", userID)
	} else if sub, exists := claims["sub"]; exists {
		uidStr = fmt.Sprintf("%v", sub)
	} else {
		return uuid.Nil, fmt.Errorf("no uid/user_id/sub found in token claims")
	}
	
	if uidStr == "" {
		return uuid.Nil, fmt.Errorf("uid is empty")
	}
	
	// Парсим в UUID
	userID, err := uuid.Parse(uidStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format: %s, error: %w", uidStr, err)
	}
	
	return userID, nil
}
