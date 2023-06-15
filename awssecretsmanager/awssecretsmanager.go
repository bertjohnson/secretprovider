// Package awssecretsmanager hosts the AWSSecretsManager type.
package awssecretsmanager

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	jsoniter "github.com/json-iterator/go"

	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
)

// AWSSecretsManager provides methods for interacting with AWS Secrets Manager.
type AWSSecretsManager struct {
	ID string

	secretsManager *secretsmanager.SecretsManager
}

var (
	// Marshaller.
	json jsoniter.API
)

// New creates a matching secret store implementation.
func New(ctx context.Context, secretStore *secretprovidertype.SecretProvider) (*AWSSecretsManager, error) {
	// Validate parameters.
	if ctx == nil {
		return nil, errors.New("context is required")
	}
	if secretStore == nil {
		return nil, errors.New("secret store configuration is required")
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, secretStore.ID) // nolint

	// Log.
	logger.Verbose(ctx, "Creating AWS Secrets Manager client.")

	// Create AWS session.
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(secretStore.ClientID, secretStore.ClientSecret, secretStore.ClientToken),
		Region:      &secretStore.Region,
	})
	if err != nil {
		return nil, err
	}

	// Create AWS Secrets Manager.
	awsSecretManagerClient := AWSSecretsManager{
		ID:             secretStore.ID,
		secretsManager: secretsmanager.New(sess),
	}

	// Log.
	logger.Verbose(ctx, "Created AWS Secrets Manager client.")

	return &awsSecretManagerClient, nil
}
