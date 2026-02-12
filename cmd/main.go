package main

import (
	"github.com/WebCraftersGH/User-service/internal/config"
	"github.com/WebCraftersGH/User-service/internal/controller"
	"github.com/WebCraftersGH/User-service/internal/repositories/user_repo"
	"github.com/WebCraftersGH/User-service/internal/usecase"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

var cfg *config.Config

func main() {
	cfg = config.Load()

	db, err := NewGORMConnection()
	if err != nil {

	}

	userREPO := userrepo.NewUserRepo(db)
	userSVC := usecase.NewUserService(userREPO)
	userCTRL := controller.NewUserController(userSVC)

	router := mux.NewRouter()
	userCTRL.RegisterRoutes(router)

	http.ListenAndServe(":"+strconv.Itoa(cfg.HTTPPort), router)
}

func NewGORMConnection() (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{})

	return db, err
}
