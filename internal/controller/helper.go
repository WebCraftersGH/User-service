package controller

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/WebCraftersGH/User-service/internal/domain"
)

func getUserToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", domain.ErrUnauthorized 
	}

	return cookie.Value, nil
}

func checkUserAuth(token string) error {
	return nil
}

func getTokenPayload(token string) (map[string]string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, domain.ErrUnauthorized
	}

	payloadB, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	var payload map[string]string
	if err := json.Unmarshal(payloadB, &payload); err != nil {
		return nil, domain.ErrUnauthorized
	}

	return payload, nil
}

func checkRights(userID string, r *http.Request) error {
	token, err := getUserToken(r)
	if err != nil {
		return err
	}

	err = checkUserAuth(token)
	if err != nil {
		return err
	}

	payload, err := getTokenPayload(token)
	if err != nil {
		return err
	}

	if userID != payload["id"] {
		return domain.ErrUnauthorized
	}
	return nil
}
