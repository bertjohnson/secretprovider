package awssecretsmanager

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
)

// ListSecrets lists secret paths.
func (a *AWSSecretsManager) ListSecrets(ctx context.Context, pathChannel chan string, errorChannel chan error) {
	// Validate parameters.
	if ctx == nil {
		errorChannel <- errors.New("context is required")

		close(pathChannel)
		close(errorChannel)

		return
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, a.ID) // nolint

	// List secrets.
	secrets, err := a.secretsManager.ListSecrets(&secretsmanager.ListSecretsInput{})
	if err != nil {
		errorChannel <- err

		close(pathChannel)
		close(errorChannel)

		return
	}
	for _, secret := range secrets.SecretList {
		pathChannel <- *secret.Name
	}

	// Log.
	logger.Verbose(ctx, "Listed secrets.")

	close(pathChannel)
	close(errorChannel)
}
