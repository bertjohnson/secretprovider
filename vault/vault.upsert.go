package vault

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	"github.com/bertjohnson/logger/types/env"
)

// UpsertSecret creates or updates a secret.
func (v *Vault) UpsertSecret(ctx context.Context, path string, data map[string]interface{}) error {
	// Validate parameters.
	if ctx == nil {
		return errors.New("context is required")
	}
	if path == "" {
		return errors.New("path is required")
	}
	if !strings.HasPrefix(path, "secret/") {
		path = "secret/" + path
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, v.ID) // nolint

	// Create secret.
	_, err := v.client.Logical().Write(path, data)
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
