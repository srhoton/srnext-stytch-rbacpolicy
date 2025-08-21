package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/srnext/stytch-rbacpolicy-lambda/internal/config"
	"github.com/srnext/stytch-rbacpolicy-lambda/internal/handler"
	"github.com/stytchauth/stytch-management-go/v2/pkg/api"
	"go.uber.org/zap"
)

func main() {
	logger, err := initLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("Starting Lambda function")

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	logger.Info("Initializing Stytch client", 
		zap.String("project_id", cfg.ProjectID),
		zap.String("workspace_key_id", cfg.WorkspaceKeyID))
	
	client := api.NewClient(cfg.WorkspaceKeyID, cfg.WorkspaceKeySecret)

	h := handler.NewHandler(client.RBACPolicy, cfg.ProjectID, logger)

	lambda.StartWithContext(context.Background(), h.HandleRequest)
}

func initLogger() (*zap.Logger, error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
