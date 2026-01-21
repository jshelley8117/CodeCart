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
	Logger              *zap.Logger
}

func NewCloudFunctionService(cfClient *client.CloudFunctionClient, helloWorldUrl string, logger *zap.Logger) CloudFunctionService {
	return CloudFunctionService{
		CloudFunctionClient: cfClient,
		HelloWorldURL:       helloWorldUrl,
		Logger:              logger.Named("cloud_function_service"),
	}
}

func (cfs CloudFunctionService) GetHelloWorld(ctx context.Context) (*client.HelloWorldResponse, error) {
	zLog := utils.FromContext(ctx, cfs.Logger)
	zLog.Debug("entered GetHelloWorld")

	response, err := cfs.CloudFunctionClient.InvokeHelloWorld(ctx, cfs.HelloWorldURL)
	if err != nil {
		zLog.Error("cloud function invocation failed", zap.Error(err))
		return nil, err
	}

	return response, nil
}
