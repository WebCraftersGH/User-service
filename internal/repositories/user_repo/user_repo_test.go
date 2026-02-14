package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/WebCraftersGH/User-service/internal/domain"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/golang-migrate/migrate/v4"
	mgpg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	testDB      *gorm.DB
	pgContainer testcontainers.Container
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	var err error
	pgContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("container start: %v", err)
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		log.Fatalf("container host: %v", err)
	}
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("container port: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=test password=test dbname=testdb sslmode=disable TimeZone=UTC",
		host, port.Port(),
	)

	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("gorm open: %v", err)
	}

	sqlDB, err := testDB.DB()
	if err != nil {
		log.Fatalf("db handle: %v", err)
	}

	migrationsPath := filepath.Join(projectRootDir(), "migrations") // поменяй на свою папку
	if err := applyMigrations(sqlDB, "testdb", migrationsPath); err != nil {
		log.Fatalf("apply migrations: %v", err)
	}

	code := m.Run()

	_ = pgContainer.Terminate(ctx)
	os.Exit(code)
}

func applyMigrations(sqlDB *sql.DB, dbName string, migrationsDir string) error {
	driver, err := mgpg.WithInstance(sqlDB, &mgpg.Config{
		DatabaseName: dbName,
	})
	if err != nil {
		return err
	}

	src := "file://" + migrationsDir

	mg, err := migrate.NewWithDatabaseInstance(src, "postgres", driver)
	if err != nil {
		return err
	}

	log.Printf("migrations dir: %s", migrationsDir)
	entries, _ := os.ReadDir(migrationsDir)
	for _, e := range entries {
		log.Println(e.Name())
	}

	err = mg.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

func projectRootDir() string {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		wd, _ := os.Getwd()
		return wd
	}
	dir := filepath.Dir(thisFile)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return filepath.Dir(thisFile)
		}
		dir = parent
	}
}

func withTx(t *testing.T) *gorm.DB {
	t.Helper()
	tx := testDB.Begin()
	if tx.Error != nil {
		t.Fatalf("begin tx: %v", tx.Error)
	}
	t.Cleanup(func() {
		_ = tx.Rollback().Error
	})
	return tx
}

func TestUserRepo_CreateUser(t *testing.T) {
	ctx := context.Background()
	tx := withTx(t) // Begin + Rollback в Cleanup

	repo := NewUserRepo(tx, nil)

	u := domain.User{
		Username: "someName",
		Sex:      domain.NewSexEnum("Male"),
	}

	// act
	userFromDB, err := repo.Create(ctx, u)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	// assert: айди проставился
	if userFromDB.ID == uuid.Nil {
		t.Fatalf("expected ID to be set")
	}

	// assert: запись реально в базе
	var got GormUser
	if err := tx.WithContext(ctx).First(&got, "id = ?", userFromDB.ID).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}

	if got.Email != userFromDB.Email.String() || got.Username != userFromDB.Username {
		t.Fatalf("mismatch: got=%+v want=%+v", got, userFromDB)
	}
}
