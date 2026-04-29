package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/WebCraftersGH/User-service/internal/adapters/kafka"
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
	logging "github.com/WebCraftersGH/User-service/pkg/logging"
)

var cfg *config.Config

func main() {
	_ = godotenv.Load()
	cfg = config.Load()

	logger, closer, err := logging.New(cfg.LoggingLevel)
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	db, err := NewGORMConnection()
	if err != nil {
		logger.WithError(err).Error("db connection error")
	}

	userREPO := userrepo.NewUserRepo(db, logger)
	userSVC := usecase.NewUserService(userREPO, logger)

	userHandler := handlers.NewUserHandler(userSVC, logger)
	docsHandler := docsH.NewDocsHandler()
	healthHandler := handlers.NewHealthHandler()

	authCl := authclient.New(cfg.AUTH_SERVICE_BASE_URL, logger)

	routerN := transporthttp.NewRouter(
		userHandler,
		healthHandler,
		docsHandler,
		authCl,
		logger,
		cfg.DEBUG_MODE,
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
	consumer, err := kafka.NewKafkaConsumer(kafkaConfig, userSVC, logger)
	if err != nil {
		logger.WithError(err).Error("kafka connection error")
	}
	go func() {
		consumer.Start()
	}()

	if err := http.ListenAndServe(":"+strconv.Itoa(cfg.HTTPPort), routerN); err != nil {
		logger.WithError(err).Error("http error")
		return
	}
}

func NewGORMConnection() (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})

	return db, err
}
