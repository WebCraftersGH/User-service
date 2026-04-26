package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/WebCraftersGH/User-service/internal/domain"
	"github.com/WebCraftersGH/User-service/internal/requestctx"
	svc "github.com/WebCraftersGH/User-service/internal/usecase"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	usecase svc.UserService
}

func NewUserHandler(usecase svc.UserService) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := requestctx.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	u, err := h.usecase.GetUser(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			writeError(w, http.StatusNotFound, "user not found")
			return
		case errors.Is(err, domain.InternalError):
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}

	userResponse := toUserResponse(u)
	writeJSON(w, http.StatusOK, userResponse)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["uuid"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "parse uuid error")
		return
	}

	u, err := h.usecase.GetUser(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			writeError(w, http.StatusNotFound, "user not found")
			return
		case errors.Is(err, domain.InternalError):
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}

	userResponse := toUserResponse(u)
	writeJSON(w, http.StatusOK, userResponse)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := requestctx.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := h.usecase.DeleteUser(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			writeError(w, http.StatusNotFound, "user not found")
			return
		case errors.Is(err, domain.InternalError):
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := requestctx.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var uUpdateReq UserUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&uUpdateReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid json input json data")
		return
	}

	domainUser := domain.User{
		ID:       userID,
		Username: uUpdateReq.Username,
		FIO:      uUpdateReq.FIO,
		BIO:      uUpdateReq.BIO,
		Sex:      domain.NewSexEnum(uUpdateReq.Sex),
		Birthday: uUpdateReq.Birthday,
	}

	u, err := h.usecase.UpdateUser(r.Context(), domainUser)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			writeError(w, http.StatusNotFound, "user not found")
			return
		case errors.Is(err, domain.InternalError):
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}

	userResponse := toUserResponse(u)
	writeJSON(w, http.StatusOK, userResponse)
}
