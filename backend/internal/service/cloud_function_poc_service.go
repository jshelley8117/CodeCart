package service

import (
	"context"

	"github.com/jshelley8117/CodeCart/internal/client"
	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type CloudFunctionConfig struct {
	HelloWorldURL  string
	HelloWorld2URL string
}

type CloudFunctionService struct {
	CloudFunctionClient *client.CloudFunctionClient
	Config              CloudFunctionConfig
}

func NewCloudFunctionService(cfClient *client.CloudFunctionClient, cfg CloudFunctionConfig) CloudFunctionService {
	return CloudFunctionService{
		CloudFunctionClient: cfClient,
		Config:              cfg,
	}
}

func (cfs CloudFunctionService) GetHelloWorld(ctx context.Context) (*model.HelloWorldResponse, error) {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("entered GetHelloWorld")

	response, err := cfs.CloudFunctionClient.InvokeHelloWorld(ctx, cfs.Config.HelloWorldURL)
	if err != nil {
		zLog.Error("cloud function invocation failed", zap.Error(err))
		return nil, err
	}

	return response, nil
}

func (cfs CloudFunctionService) GetHelloWorld2(ctx context.Context) (*model.HelloWorldResponse, error) {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("entered GetHelloWorld2")

	response, err := cfs.CloudFunctionClient.InvokeHelloWorld2(ctx, cfs.Config.HelloWorld2URL)
	if err != nil {
		zLog.Error("cloud function invocation failed", zap.Error(err))
		return nil, err
	}

	return response, nil
}
