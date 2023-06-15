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
)

// UpsertSecret creates or updates a secret.
func (a *AWSSecretsManager) UpsertSecret(ctx context.Context, path string, data map[string]interface{}) error {
	// Validate parameters.
	if ctx == nil {
		return errors.New("context is required")
	}
	if path == "" {
		return errors.New("path is required")
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, a.ID) // nolint

	// Create secret.
	dataBytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	_, err = a.secretsManager.PutSecretValue(&secretsmanager.PutSecretValueInput{
		SecretBinary: dataBytes,
		SecretId:     &path,
	})
	// If the secret does not exist, create it.
	if err != nil && strings.HasPrefix(err.Error(), "ResourceNotFoundException") {
		_, err = a.secretsManager.CreateSecret(&secretsmanager.CreateSecretInput{
			Name:         &path,
			SecretBinary: dataBytes,
		})
	}
	if err != nil {
		return err
	}

	// Log.
	if os.Getenv(env.Debug) != "" {
		logger.Verbose(ctx, "Upserted secret: "+path)
	} else {
		logger.Info(ctx, "Upserted secret.")
	}

	return nil
}
