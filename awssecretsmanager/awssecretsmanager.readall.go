package awssecretsmanager

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
)

// ReadAllSecrets reads all secrets.
func (a *AWSSecretsManager) ReadAllSecrets(ctx context.Context, secretChannel chan *secretprovidertype.Secret, errorChannel chan error) {
	// Validate parameters.
	if ctx == nil {
		errorChannel <- errors.New("context is required")

		close(secretChannel)
		close(errorChannel)

		return
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, a.ID) // nolint

	// Read all secrets.
	secrets, err := a.secretsManager.ListSecrets(&secretsmanager.ListSecretsInput{})
	if err != nil {
		errorChannel <- err

		close(secretChannel)
		close(errorChannel)

		return
	}
	for _, secretValue := range secrets.SecretList {
		secret, err := a.ReadSecret(ctx, *secretValue.Name)
		if err != nil {
			errorChannel <- err

			close(secretChannel)
			close(errorChannel)

			return
		}

		// Return secret.
		secretChannel <- secret
	}

	// Log.
	logger.Verbose(ctx, "Read all secrets.")

	close(secretChannel)
	close(errorChannel)
}
