package utility_test

import (
	"context"
	"file-manager/internal/configuration"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func RunContainer(t *testing.T, pathMigration string) configuration.DatasourceConfig {
	ctx := context.Background()
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:15.3-alpine",
		postgres.WithDatabase("file-manager"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)))

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	endpoint, err := pgContainer.PortEndpoint(ctx, "5432/tcp", "")

	if err != nil {
		t.Fatalf("failed get endpoint for container, pgcontainer: %s", err)
	}

	host := strings.Split(endpoint, ":")[0]
	port, err := strconv.Atoi(strings.Split(endpoint, ":")[1])
	if err != nil {
		t.Fatalf("failed get port for container, error: %s", err)
	}

	datasourceConfig := configuration.DatasourceConfig{
		Host:     host,
		Port:     int32(port),
		Username: "postgres",
		Password: "postgres",
		Driver:   "postgres",
		Database: "file-manager",
	}

	if err := datasourceConfig.Migration(pathMigration); err != nil {
		t.Fatalf("failed migration up with error: %s", err)
	}

	if err := configuration.RunDatabase(datasourceConfig); err != nil {
		t.Fatalf("failed start database, error: %s", err)
	}
	return datasourceConfig
}
