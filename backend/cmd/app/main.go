package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/jshelley8117/CodeCart/internal/middleware"
	"github.com/jshelley8117/CodeCart/internal/resource"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/api/impersonate"
)

const EXIT_STATUS = 1

type ResourceConfig struct {
	GCloudDB *sql.DB
	Logger   *zap.Logger
}

func main() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Panic("cannot load env vars")
	}

	logger, err := utils.NewLogger(utils.Config{Env: os.Getenv("ENV")})
	if err != nil {
		log.Panic("cannot instantiate logger")
	}

	targetSA := os.Getenv("GCP_IMP_SA")
	if targetSA == "" {
		logger.Error("GCP_IMP_SA is not set")
		os.Exit(EXIT_STATUS)
	}

	ctx := context.Background()
	ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: targetSA,
		Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
	})
	if err != nil {
		logger.Error("failed to create impersonated token source: %v", zap.Error(err))
		os.Exit(EXIT_STATUS)
	}

	tok, err := ts.Token()
	if err != nil {
		logger.Error("failed to mint impersonated token: %v", zap.Error(err))
		os.Exit(EXIT_STATUS)
	}

	logger.Debug("impersonation OK; token expires at %s (in ~%s)", zap.String("expiry", tok.Expiry.Format(time.RFC3339)), zap.String("duration", time.Until(tok.Expiry).Round(time.Second).String()))

	reusableTS := oauth2.ReuseTokenSource(tok, ts)

	dbHandle, err := resource.NewGCloudDB(reusableTS)
	if err != nil {
		logger.Error("could not connect to db: %v", zap.Error(err))
		os.Exit(EXIT_STATUS)
	}

	logger.Debug("connected to db")

	mux := http.NewServeMux()
	SetupRoutes(mux, ResourceConfig{
		GCloudDB: dbHandle,
		Logger:   logger,
	})

	handler := middleware.Recoverer(logger)(middleware.RequestLogger(logger)(mux))

	server := http.Server{
		Addr:    ":8081",
		Handler: handler,
	}
	logger.Debug("go server initiating", zap.String("addr", server.Addr))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("error starting server", zap.Error(err))
	}
}
