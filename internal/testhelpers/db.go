package testhelpers

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/kiennyo/syncwatch-be/internal/config"
	"github.com/kiennyo/syncwatch-be/internal/db"
)

type TestingDB struct {
	DB *pgxpool.Pool
}

var cached *TestingDB

//nolint:revive,cognitive-complexity
func CreateTestDB(ctx context.Context) (*TestingDB, error) {
	if cached != nil {
		slog.Info("Getting cached instance")
		return cached, nil
	}

	name := fmt.Sprintf("syncwatch-integration-tests-db-%s", uuid.New().String())
	user := uuid.New().String()
	password := uuid.New().String()

	container, err := postgres.RunContainer(
		ctx,
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Name: "syncwatch-integration-tests-db",
			},
		}),
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		postgres.WithDatabase(name),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		slog.Error("Failed to create postgres container")
		return nil, err
	}

	// set DB_URL for migrations
	ip, _ := container.ContainerIP(ctx)
	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, ip, name)

	err = os.Setenv("DB_URL", connectionString)
	if err != nil {
		return nil, err
	}

	// run migrations
	cmd := exec.Command("task", "db:run-migration")
	stderr, _ := cmd.StderrPipe()

	if err = cmd.Start(); err != nil {
		log.Printf("cmd.Start: %v", err)
	}

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)

	// for debugging
	// for scanner.Scan() {
	//	m := scanner.Text()
	//	fmt.Println(m)
	// }

	if err = cmd.Wait(); err != nil {
		return nil, err
	}

	s, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		slog.Error("Failed getting connection string")
		return nil, err
	}

	pool, err := db.New(ctx, config.DB{
		URL:         s,
		MaxOpenConn: 1,
		MaxIdleConn: 1,
		MaxIdleTime: "1m",
	})
	if err != nil {
		slog.Error("Failed connecting to database")
		return nil, err
	}

	cached = &TestingDB{DB: pool}

	return cached, nil
}
