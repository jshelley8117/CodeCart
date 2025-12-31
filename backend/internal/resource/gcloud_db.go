package resource

import (
	"database/sql"
	"fmt"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv5"
	"golang.org/x/oauth2"
)

func NewGCloudDB(tokenSource oauth2.TokenSource) (*sql.DB, error) {
	dbUser, err := getEnv("GCP_DB_USER")
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

	_, err = pgxv5.RegisterDriver("cloudsqlpostgres",
		cloudsqlconn.WithIAMAuthN(),
		cloudsqlconn.WithIAMAuthNTokenSources(tokenSource, tokenSource),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register driver: %w", err)
	}

	dsn := fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable",
		dbUser,
		dbName,
		instanceName,
	)

	db, err := sql.Open("cloudsqlpostgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
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
