package vault

import (
	"context"
	"errors"

	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
)

// ListSecrets lists secret paths.
func (v *Vault) ListSecrets(ctx context.Context, pathChannel chan string, errorChannel chan error) {
	// Validate parameters.
	if ctx == nil {
		errorChannel <- errors.New("context is required")

		close(pathChannel)
		close(errorChannel)

		return
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, v.ID) // nolint

	// List secrets.
	// TODO.

	// Log.
	logger.Verbose(ctx, "Listed secrets.")

	close(pathChannel)
	close(errorChannel)
}
