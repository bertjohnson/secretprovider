package awssecretsmanager

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	"github.com/bertjohnson/logger/types/env"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
)

// ReadSecret returns a secret.
func (a *AWSSecretsManager) ReadSecret(ctx context.Context, path string) (secret *secretprovidertype.Secret, err error) {
	// Validate parameters.
	if ctx == nil {
		return nil, errors.New("context is required")
	}
	if path == "" {
		return nil, errors.New("path is required")
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, a.ID) // nolint

	// Read secret.
	secretValue, err := a.secretsManager.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: &path,
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "ResourceNotFoundException") {
			return nil, errors.New("not found")
		}

		return nil, err
	}
	if secretValue == nil {
		if os.Getenv(env.Debug) != "" {
			logger.Verbose(ctx, "Secret does not exist: "+path)
		} else {
			logger.Verbose(ctx, "Secret does not exist.")
		}

		return nil, errors.New("not found")
	}

	// Check if the secret is binary or a string.
	secret = new(secretprovidertype.Secret)
	if secretValue.SecretString != nil && len(*secretValue.SecretString) > 0 {
		err = json.Unmarshal([]byte(*secretValue.SecretString), &secret.Data)
	} else {
		err = json.Unmarshal(secretValue.SecretBinary, &secret.Data)
	}
	if err != nil {
		return nil, err
	}
	secret.Path = path

	// Log.
	if os.Getenv(env.Debug) != "" {
		logger.Verbose(ctx, "Read secret: "+path)
	} else {
		logger.Verbose(ctx, "Read secret.")
	}

	return secret, nil
}
