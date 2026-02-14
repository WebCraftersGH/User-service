package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/WebCraftersGH/User-service/internal/domain"
	"github.com/WebCraftersGH/User-service/internal/usecase"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type userController struct {
	svc usecase.UserService
	lg  usecase.Logger
}

func NewUserController(svc usecase.UserService, lg usecase.Logger) *userController {
	return &userController{svc: svc, lg: lg}
}

func (c *userController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/users/me", c.GetMe).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/users/{uuid}", c.GetUserByID).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/users/{uuid}", c.DeleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/api/v1/users/{uuid}", c.UpdateUser).Methods(http.MethodPut)
}

func (c *userController) GetMe(w http.ResponseWriter, r *http.Request) {
	c.lg.Debug("[User-controller][GetMe][DEBUG] - Start function")
	token, err := getUserToken(r)
	if err != nil {

	}

	if err := checkUserAuth(token); err != nil {

	}

	payload, err := getTokenPayload(token)
	if err != nil {

	}

	id, ok := payload["id"]
	if !ok {

	}

	uID, err := uuid.Parse(id)
	if err != nil {

	}

	u, err := c.svc.GetUser(r.Context(), uID)
	if err != nil {

	}

	uR := toUserResponse(u)

	c.lg.Debug("[User-controller][GetMe][DEBUG] - End function")
	json.NewEncoder(w).Encode(uR)
}

func (c *userController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	c.lg.Debug("[User-controller][GetUserByID][DEBUG] - Start function")

	id, err := uuid.Parse(mux.Vars(r)["uuid"])
	if err != nil {
		http.Error(w, "parse user uuid error", http.StatusBadRequest)
	}

	u, err := c.svc.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, domain.ErrUserNotFound.Error(), http.StatusNotFound)
		}
		http.Error(w, domain.InternalError.Error(), http.StatusInternalServerError)
	}

	uR := toUserResponse(u)

	c.lg.Debug("[User-controller][GetUserByID][DEBUG] - End function")
	json.NewEncoder(w).Encode(uR)
}

func (c *userController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	c.lg.Debug("[User-controller][DeleteUser][DEBUG] - Start function")

	id, err := uuid.Parse(mux.Vars(r)["uuid"])
	if err != nil {
		http.Error(w, "parse user uuid error", http.StatusBadRequest)
		return
	}

	if err := checkRights(id.String(), r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := c.svc.DeleteUser(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.lg.Debug("[User-controller][GetUserByID][DEBUG] - End function")
	w.WriteHeader(http.StatusOK)
}

func (c *userController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	c.lg.Debug("[User-controller][UpdateUser][DEBUG] - Start function")

	id, err := uuid.Parse(mux.Vars(r)["uuid"])
	if err != nil {
		http.Error(w, "parse user uuid error", http.StatusBadRequest)
		return
	}

	if err := checkRights(id.String(), r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var uUpdateReq UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&uUpdateReq); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	domainUser := domain.User{
		ID:       id,
		Username: uUpdateReq.Username,
		FIO:      uUpdateReq.FIO,
		BIO:      uUpdateReq.BIO,
		Sex:      domain.NewSexEnum(uUpdateReq.Sex),
		Birthday: uUpdateReq.Birthday,
	}

	u, err := c.svc.UpdateUser(r.Context(), domainUser)
	if err != nil {

	}

	uR := toUserResponse(u)

	c.lg.Debug("[User-controller][UpdateUser][DEBUG] - End function")
	json.NewEncoder(w).Encode(uR)
}
