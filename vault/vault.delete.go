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

// DeleteSecret deletes a secret.
func (v *Vault) DeleteSecret(ctx context.Context, path string) error {
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

	// Delete secret.
	_, err := v.client.Logical().Delete(path)
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
