package vault

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	"github.com/bertjohnson/logger/types/env"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
)

// ReadSecret returns a secret.
func (v *Vault) ReadSecret(ctx context.Context, path string) (secret *secretprovidertype.Secret, err error) {
	// Validate parameters.
	if ctx == nil {
		return nil, errors.New("context is required")
	}
	if path == "" {
		return nil, errors.New("path is required")
	}
	if !strings.HasPrefix(path, "secret/") {
		path = "secret/" + path
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, v.ID) // nolint

	// Read secret.
	vaultSecret, err := v.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}
	if vaultSecret == nil {
		if os.Getenv(env.Debug) != "" {
			logger.Verbose(ctx, "Secret does not exist: "+path)
		} else {
			logger.Verbose(ctx, "Secret does not exist.")
		}

		return nil, errors.New("not found")
	}
	secret = new(secretprovidertype.Secret)
	secret.Data = vaultSecret.Data
	secret.Path = path

	// Log.
	if os.Getenv(env.Debug) != "" {
		logger.Verbose(ctx, "Read secret: "+path)
	} else {
		logger.Verbose(ctx, "Read secret.")
	}

	return secret, nil
}
