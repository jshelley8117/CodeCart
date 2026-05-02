package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/api/impersonate"
)

type CloudFunctionClient struct {
	HttpClient          *http.Client
	ServiceAccountEmail string
	mu                  sync.Mutex
	tokenSources        map[string]oauth2.TokenSource // mutex protected map
}

func NewCloudFunctionClient(serviceAccountEmail string) *CloudFunctionClient {
	return &CloudFunctionClient{
		HttpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		ServiceAccountEmail: serviceAccountEmail,
		tokenSources:        make(map[string]oauth2.TokenSource),
	}
}

func invokeFunction[T any](ctx context.Context, c *CloudFunctionClient, url, method string, requestBody any) (*T, error) {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("invoking cloud function",
		zap.String("url", url),
		zap.String("method", method))

	var req *http.Request
	var err error

	if requestBody != nil {
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			z.Error("failed to marshal request body", zap.Error(err))
			return nil, err
		}

		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyBytes))
		if err != nil {
			z.Error("failed to create request", zap.Error(err))
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			z.Error("failed to create request", zap.Error(err))
			return nil, err
		}
	}

	idToken, err := c.getIdToken(ctx, url)
	if err != nil {
		z.Error("failed to get ID token", zap.Error(err))
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", idToken))

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		z.Error("failed to invoke cloud function", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		z.Error("failed to read response body", zap.Error(err))
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		z.Error("cloud function returned a non-success status",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(body)))
		return nil, err
	}

	var response T
	if err := json.Unmarshal(body, &response); err != nil {
		z.Error("failed to unmarshal response", zap.Error(err))
		return nil, err
	}
	return &response, nil
}

func (c *CloudFunctionClient) getIdToken(ctx context.Context, audience string) (string, error) {
	// locking mutex to protect tokenSources map from concurrent writes during initialization
	c.mu.Lock()
	ts, ok := c.tokenSources[audience]
	if !ok {
		rawTokenSource, err := impersonate.IDTokenSource(ctx, impersonate.IDTokenConfig{
			Audience:        audience,
			TargetPrincipal: c.ServiceAccountEmail,
			IncludeEmail:    true,
		})
		if err != nil {
			c.mu.Unlock()
			return "", err
		}
		ts = oauth2.ReuseTokenSource(nil, rawTokenSource)
		c.tokenSources[audience] = ts
	}
	c.mu.Unlock()

	t, err := ts.Token()
	if err != nil {
		return "", err
	}
	return t.AccessToken, nil
}

// --- GCP INVOKE FUNCTIONS DEFINED BELOW

func (c *CloudFunctionClient) InvokeHelloWorld(ctx context.Context, url string) (*model.HelloWorldResponse, error) {
	return invokeFunction[model.HelloWorldResponse](ctx, c, url, http.MethodGet, nil)
}

func (c *CloudFunctionClient) InvokeHelloWorld2(ctx context.Context, url string) (*model.HelloWorldResponse, error) {
	return invokeFunction[model.HelloWorldResponse](ctx, c, url, http.MethodGet, nil)
}
