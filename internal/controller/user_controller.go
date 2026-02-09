package controller

import (
	"github.com/WebCraftersGH/User-service/internal/usecase"
	"github.com/gorilla/mux"
	"net/http"
)

type userController struct {
	svc usecase.UserService
}

func NewUserController(svc usecase.UserService) *userController {
	return &userController{svc: svc}
}

func (c *userController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/users/me", c.GetMe).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/users/{id}", c.GetUserByID).Methods(http.MethodGet)
}

func (c *userController) GetMe(w http.ResponseWriter, r *http.Request)       {}
func (c *userController) GetUserByID(w http.ResponseWriter, r *http.Request) {}
