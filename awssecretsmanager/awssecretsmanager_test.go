package awssecretsmanager

import (
	"context"
	"os"
	"testing"

	"github.com/bertjohnson/logger"
	"github.com/bertjohnson/logger/types/env"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
	"github.com/bertjohnson/startup"
)

var (
	// Context.
	ctx context.Context

	// AWS Secrets Manager client.
	awsSecretsManager *AWSSecretsManager
)

// TestMain runs tests.
func TestMain(m *testing.M) {
	// Wait for configuration to initialize.
	ctx = context.Background()
	go func() {
		startup.Ready()
	}()
	startup.Wait(ctx, startup.PackageType) // Declare that the configuration is ready.

	// Wait for logger.
	logger.Wait(ctx)

	// Create AWS Secrets Manager client.
	var err error
	awsSecretsManager, err = New(ctx, &secretprovidertype.SecretProvider{
		ID:           "5587e3a7-d5cf-44de-b403-8aa01a93670c",
		ClientID:     os.Getenv(env.SecretProviderClientID),
		ClientSecret: os.Getenv(env.SecretProviderClientSecret),
		Region:       os.Getenv(env.SecretProviderRegion),
		Type:         "AWSSecretsManager",
	})
	if err != nil {
		logger.Fatal(ctx, "Error creating AWS Secrets Manager client: "+err.Error())
	}

	// Run tests.
	retCode := m.Run()

	// Exit.
	os.Exit(retCode)
}
