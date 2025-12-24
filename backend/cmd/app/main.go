package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/jshelley8117/CodeCart/internal/resource"
	"golang.org/x/oauth2"
	"google.golang.org/api/impersonate"
)

func main() {
	if err := godotenv.Load("./../../.env"); err != nil {
		log.Panic("cannot load env vars")
	}

	targetSA := os.Getenv("GCP_IMP_SA")
	if targetSA == "" {
		log.Panic("GCP_IMP_SA is not set")
	}

	ctx := context.Background()
	ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: targetSA,
		Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
	})
	if err != nil {
		log.Panicf("failed to create impersonated token source: %v", err)
	}

	tok, err := ts.Token()
	if err != nil {
		log.Panicf("failed to mint impersonated token: %v", err)
	}

	log.Printf("impersonation OK; token expires at %s (in ~%s)", tok.Expiry.Format(time.RFC3339), time.Until(tok.Expiry).Round(time.Second))

	_ = oauth2.ReuseTokenSource(tok, ts)

	_, err = resource.ConnectWithConnector()
	if err != nil {
		fmt.Printf("could not connect to db: %v", err)
		log.Panic("exiting main")
	}

	log.Println("connected to db")

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	log.Println("Go Server Listening on Port 8081")
	server.ListenAndServe()
}
