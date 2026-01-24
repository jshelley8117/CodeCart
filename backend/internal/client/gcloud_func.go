package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/api/impersonate"
)

type CloudFunctionClient struct {
	HttpClient          *http.Client
	TokenSource         oauth2.TokenSource
	ServiceAccountEmail string
}

func NewCloudFunctionClient(tokenSource oauth2.TokenSource, serviceAccountEmail string) *CloudFunctionClient {
	return &CloudFunctionClient{
		HttpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		TokenSource:         tokenSource,
		ServiceAccountEmail: serviceAccountEmail,
	}
}

func (cfc *CloudFunctionClient) invokeFunction(ctx context.Context, url, method string, requestBody, response any) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("invoking cloud function",
		zap.String("url", url),
		zap.String("method", method))

	var req *http.Request
	var err error

	if requestBody != nil {
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			zLog.Error("failed to marshal request body", zap.Error(err))
			return fmt.Errorf("failed to marshal request body: %w", err)
		}

		req, err = http.NewRequestWithContext(ctx, method, url, io.NopCloser(bytes.NewReader(bodyBytes)))
		if err != nil {
			zLog.Error("failed to create request", zap.Error(err))
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			zLog.Error("failed to create request", zap.Error(err))
			return fmt.Errorf("failed to create request: %w", err)
		}
	}

	idToken, err := cfc.getIdToken(ctx, url)
	if err != nil {
		zLog.Error("failed to get ID token", zap.Error(err))
		return fmt.Errorf("failed to get ID token: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", idToken))

	resp, err := cfc.HttpClient.Do(req)
	if err != nil {
		zLog.Error("failed to invoke cloud function", zap.Error(err))
		return fmt.Errorf("failed to invoke cloud function: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zLog.Error("failed to read response body", zap.Error(err))
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		zLog.Error("cloud function returned a non-success status",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(body)))
		return fmt.Errorf("cloud function returned status %d: %s", resp.StatusCode, string(body))
	}

	if response != nil {
		if err := json.Unmarshal(body, response); err != nil {
			zLog.Error("failed to unmarshal response", zap.Error(err))
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

func (cfc *CloudFunctionClient) getIdToken(ctx context.Context, audience string) (string, error) {
	ts, err := impersonate.IDTokenSource(ctx, impersonate.IDTokenConfig{
		Audience:        audience,
		TargetPrincipal: cfc.ServiceAccountEmail,
		IncludeEmail:    true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create impersonated ID token source: %w", err)
	}

	token, err := ts.Token()
	if err != nil {
		return "", fmt.Errorf("failed to get impersonated ID token: %w", err)
	}

	return token.AccessToken, nil
}

func (cfc *CloudFunctionClient) InvokeHelloWorld(ctx context.Context, url string) (*HelloWorldResponse, error) {
	var response HelloWorldResponse
	if err := cfc.invokeFunction(ctx, url, http.MethodGet, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type HelloWorldResponse struct {
	Message string `json:"message"`
}
