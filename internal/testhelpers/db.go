package testhelpers

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
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

func CreateTestDB(ctx context.Context) (*TestingDB, error) {
	if cached != nil {
		fmt.Println("Getting cached instance")
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
		log.Fatalf("cmd.Start: %v", err)
	}

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	if err = cmd.Wait(); err != nil {
		var exiterr *exec.ExitError
		if errors.As(err, &exiterr) {
			fmt.Println(fmt.Sprintf("Exit Status: %d", exiterr.ExitCode()))
		}
	}

	s, err := container.ConnectionString(ctx, "sslmode=disable")

	pool, err := db.New(ctx, config.DB{
		URL:         s,
		MaxOpenConn: 1,
		MaxIdleConn: 1,
		MaxIdleTime: "1m",
	})

	cached = &TestingDB{DB: pool}

	return cached, nil
}
