package resource

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func NewGCloudDBHandle() (*sql.DB, error) {
	dbUser, err := getEnv("GCP_IMP_SA")
	if err != nil {
		return nil, err
	}
	dbPassword, err := getEnv("GCP_DB_PASS")
	if err != nil {
		return nil, err
	}
	dbName, err := getEnv("GCP_DB_NAME")
	if err != nil {
		return nil, err
	}
	instanceName, err := getEnv("GCP_INSTANCE_CONNECTION_NAME")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("user=%s password=%s database=%s", dbUser, dbPassword, dbName)
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	dialer, err := createDialer()
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloud SQL dialer: %w", err)
	}

	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return dialer.Dial(ctx, instanceName)
	}

	dbURI := stdlib.RegisterConnConfig(config)
	db, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	return db, nil
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("required environment variable %s is not set", key)
	}
	return value, nil
}

func createDialer() (*cloudsqlconn.Dialer, error) {
	opts := []cloudsqlconn.Option{cloudsqlconn.WithLazyRefresh()}

	if os.Getenv("PRIVATE_IP") != "" {
		opts = append(opts, cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()))
	}

	return cloudsqlconn.NewDialer(context.Background(), opts...)
}
