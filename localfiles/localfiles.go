// Package localfiles hosts the LocalFiles type.
package localfiles

import (
	"context"
	"errors"
	"os"

	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	"github.com/bertjohnson/logger/types/env"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
	utilio "github.com/bertjohnson/util/io"
)

// LocalFiles provides methods for interacting with LocalFiles.
type LocalFiles struct {
	basePath string
	ID       string
}

var (
	// Path separator.
	pathSeparator = string(os.PathSeparator)
)

// New creates a matching secret store implementation.
func New(ctx context.Context, secretStore *secretprovidertype.SecretProvider) (*LocalFiles, error) {
	// Validate parameters.
	if ctx == nil {
		return nil, errors.New("context is required")
	}
	if secretStore == nil {
		return nil, errors.New("secret store configuration is required")
	}
	if secretStore.URI == "" {
		secretStore.URI = os.Getenv(env.SecretProviderURI)
		if secretStore.URI == "" {
			return nil, errors.New("secret store URI is required")
		}
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, secretStore.ID) // nolint

	// Log.
	logger.Verbose(ctx, "Creating LocalFiles client.")

	// Initialize LocalFiles client.
	localfilesClient := LocalFiles{
		ID:       secretStore.ID,
		basePath: secretStore.URI,
	}

	// Ensure the directory exists.
	err := utilio.EnsureDirectory(secretStore.URI)
	if err != nil {
		return nil, err
	}

	// Log.
	logger.Verbose(ctx, "Created LocalFiles client.")

	return &localfilesClient, nil
}
