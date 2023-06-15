// Package localfiles hosts the LocalFiles type.
package localfiles

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	utilio "github.com/bertjohnson/util/io"
)

// DeleteSecret deletes a secret.
func (l *LocalFiles) DeleteSecret(ctx context.Context, path string) error {
	// Validate parameters.
	if ctx == nil {
		return errors.New("context is required")
	}
	if path == "" {
		return errors.New("path is required")
	}
	if strings.Contains(path, "..") {
		return errors.New("path contains fobidden sequence (..): " + path)
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, l.ID) // nolint

	// Delete item.
	originalURI := l.basePath + pathSeparator + path + ".secret"
	if utilio.FileExists(originalURI) {
		err := os.Remove(originalURI)
		if err != nil {
			return err
		}
	}

	// Log.
	logger.Info(ctx, "Deleted secret.")

	return nil
}
