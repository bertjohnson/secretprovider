// Package localfiles hosts the LocalFiles type.
package localfiles

import (
	"context"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
	utilio "github.com/bertjohnson/util/io"
	json "github.com/json-iterator/go"
)

// ReadSecret returns a secret.
func (l *LocalFiles) ReadSecret(ctx context.Context, path string) (secret *secretprovidertype.Secret, err error) {
	// Validate parameters.
	if ctx == nil {
		return nil, errors.New("context is required")
	}
	if path == "" {
		return nil, errors.New("path is required")
	}
	if strings.Contains(path, "..") {
		return nil, errors.New("path contains fobidden sequence (..): " + path)
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, l.ID) // nolint

	// Read and deserialize file.
	if !strings.HasSuffix(path, ".secret") {
		path += ".secret"
	}
	data, err := ioutil.ReadFile(l.basePath + pathSeparator + utilio.NormalizePathSeparators(path))
	if err != nil {
		return nil, err
	}
	dataMap := make(map[string]interface{})
	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		return nil, err
	}

	// Return secret.
	secret = &secretprovidertype.Secret{
		Data: dataMap,
		Path: path,
	}

	// Log.
	logger.Verbose(ctx, "Read secret.")

	return secret, nil
}
