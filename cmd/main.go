package main

import (
	"github.com/WebCraftersGH/User-service/internal/adapters/kafka"
	"github.com/WebCraftersGH/User-service/internal/adapters/logging"
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
	"time"
)

var cfg *config.Config

func main() {
	_ = godotenv.Load()
	cfg = config.Load()

	lg := logging.NewLogger(cfg.LoggingLevel)

	db, err := NewGORMConnection()
	if err != nil {
		lg.Error("[MAIN][Gorm-conection][ERROR] - Gorm connection error", "gorm_err", err)
	}

	userREPO := userrepo.NewUserRepo(db, lg)
	userSVC := usecase.NewUserService(userREPO, lg)
	userCTRL := controller.NewUserController(userSVC, lg)

	router := mux.NewRouter()
	userCTRL.RegisterRoutes(router)

	kafkaConfig := &kafka.Config{
		TimeoutMS:          cfg.KafkaTimeoutMS,
		Topic:              cfg.KafkaTopic,
		GroupID:            cfg.KafkaGroupID,
		AutoOffsetStore:    cfg.KafkaAutoOffsetStore,
		ReadMessageTimeout: time.Duration(cfg.KafkaReadMessageTimeout) * time.Millisecond,
		AutoCommit:         cfg.KafkaAutoCommit,
		AutoCommitInterval: cfg.KafkaAutoCommitInterval,
		BootStrapServers:   cfg.KafkaBrokers,
	}
	consumer, err := kafka.NewKafkaConsumer(kafkaConfig, userSVC, lg)
	if err != nil {
		lg.Error("[MAIN][Kafka-conection][ERROR] - Kafka connection error", "kafka_err", err)
	}
	go func() {
		consumer.Start()
	}()

	http.ListenAndServe(":"+strconv.Itoa(cfg.HTTPPort), router)
}

func NewGORMConnection() (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})

	return db, err
}
