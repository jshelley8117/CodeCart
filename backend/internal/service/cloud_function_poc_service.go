package service

import (
	"context"

	"github.com/jshelley8117/CodeCart/internal/client"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type CloudFunctionService struct {
	CloudFunctionClient *client.CloudFunctionClient
	HelloWorldURL       string
}

func NewCloudFunctionService(cfClient *client.CloudFunctionClient, helloWorldUrl string) CloudFunctionService {
	return CloudFunctionService{
		CloudFunctionClient: cfClient,
		HelloWorldURL:       helloWorldUrl,
	}
}

func (cfs CloudFunctionService) GetHelloWorld(ctx context.Context) (*client.HelloWorldResponse, error) {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("entered GetHelloWorld")

	response, err := cfs.CloudFunctionClient.InvokeHelloWorld(ctx, cfs.HelloWorldURL)
	if err != nil {
		zLog.Error("cloud function invocation failed", zap.Error(err))
		return nil, err
	}

	return response, nil
}
