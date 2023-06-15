package localfiles

import (
	"context"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	utilio "github.com/bertjohnson/util/io"
	json "github.com/json-iterator/go"
)

// UpsertSecret creates or updates a secret.
func (l *LocalFiles) UpsertSecret(ctx context.Context, path string, data map[string]interface{}) error {
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

	// Create secret.
	if !strings.HasSuffix(path, ".secret") {
		path += ".secret"
	}
	uri := l.basePath + pathSeparator + utilio.NormalizePathSeparators(path)
	lastSlash := strings.LastIndex(uri, pathSeparator)
	uriPath := uri[0:lastSlash]
	err := utilio.EnsureDirectory(uriPath)
	if err != nil {
		return err
	}

	// Serialize data.
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write secret to file.
	err = ioutil.WriteFile(uri, dataBytes, 0644)
	if err != nil {
		return err
	}

	// Log.
	logger.Info(ctx, "Upserted secret.")

	return nil
}
