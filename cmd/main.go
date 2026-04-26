package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/WebCraftersGH/User-service/internal/adapters/kafka"
	"github.com/WebCraftersGH/User-service/internal/adapters/logging"
	"github.com/WebCraftersGH/User-service/internal/config"
	"github.com/WebCraftersGH/User-service/internal/repositories/user_repo"
	transporthttp "github.com/WebCraftersGH/User-service/internal/transport"
	"github.com/WebCraftersGH/User-service/internal/usecase"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	authclient "github.com/WebCraftersGH/User-service/internal/authclient"
	docsH "github.com/WebCraftersGH/User-service/internal/transport/http/docs"
	handlers "github.com/WebCraftersGH/User-service/internal/transport/http/handlers"
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

	userHandler := handlers.NewUserHandler(userSVC)
	docsHandler := docsH.NewDocsHandler()
	healthHandler := handlers.NewHealthHandler()

	authCl := authclient.New(cfg.AUTH_SERVICE_BASE_URL)

	routerN := transporthttp.NewRouter(
		userHandler,
		healthHandler,
		docsHandler,
		authCl,
	)

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

	http.ListenAndServe(":"+strconv.Itoa(cfg.HTTPPort), routerN)
}

func NewGORMConnection() (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})

	return db, err
}
