package vault

import (
	"context"
	"errors"

	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
)

// ReadAllSecrets reads all secrets.
func (v *Vault) ReadAllSecrets(ctx context.Context, secretChannel chan *secretprovidertype.Secret, errorChannel chan error) {
	// Validate parameters.
	if ctx == nil {
		errorChannel <- errors.New("context is required")

		close(secretChannel)
		close(errorChannel)

		return
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, v.ID) // nolint

	// Read all secrets.
	// TODO.

	// Log.
	logger.Verbose(ctx, "Read all secrets.")

	close(secretChannel)
	close(errorChannel)
}
