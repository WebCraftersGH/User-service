package main

import (
	"github.com/WebCraftersGH/User-service/internal/adapters/kafka"
	"github.com/WebCraftersGH/User-service/internal/config"
	"github.com/WebCraftersGH/User-service/internal/controller"
	"github.com/WebCraftersGH/User-service/internal/repositories/user_repo"
	"github.com/WebCraftersGH/User-service/internal/usecase"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

var cfg *config.Config

func main() {
	_ = godotenv.Load()
	cfg = config.Load()

	db, err := NewGORMConnection()
	if err != nil {

	}

	userREPO := userrepo.NewUserRepo(db)
	userSVC := usecase.NewUserService(userREPO)
	userCTRL := controller.NewUserController(userSVC)

	router := mux.NewRouter()
	userCTRL.RegisterRoutes(router)

	kafkaConfig := &kafka.Config{}
	consumer, err := kafka.NewKafkaConsumer(kafkaConfig, userSVC)
	if err != nil {

	}
	go func() {
		consumer.Start()
	}()

	http.ListenAndServe(":"+strconv.Itoa(cfg.HTTPPort), router)
}

func NewGORMConnection() (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{})

	return db, err
}
