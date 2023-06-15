package awssecretsmanager

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	"github.com/bertjohnson/logger/types/env"
)

// DeleteSecret deletes a secret.
func (a *AWSSecretsManager) DeleteSecret(ctx context.Context, path string) error {
	// Validate parameters.
	if ctx == nil {
		return errors.New("context is required")
	}
	if path == "" {
		return errors.New("path is required")
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, a.ID) // nolint

	// Delete secret.
	_, err := a.secretsManager.DeleteSecret(&secretsmanager.DeleteSecretInput{
		SecretId: &path,
	})
	if err != nil {
		return err
	}

	// Log.
	if os.Getenv(env.Debug) != "" {
		logger.Verbose(ctx, "Deleted secret: "+path)
	} else {
		logger.Info(ctx, "Deleted secret.")
	}

	return nil
}
