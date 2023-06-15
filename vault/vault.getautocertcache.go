package vault

import (
	"context"
	"encoding/base64"
	"errors"
	"os"
	"strings"

	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	"github.com/bertjohnson/logger/types/env"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
	vault "github.com/hashicorp/vault/api"
	"golang.org/x/crypto/acme/autocert"
)

// AutoCertCache implements AutoCertCache using Vault.
type AutoCertCache struct {
	client *vault.Client
	ID     string
}

// GetAutoCertCache returns an autocert-compatible cache.
func (v *Vault) GetAutoCertCache(ctx context.Context) secretprovidertype.AutoCertCache {
	return AutoCertCache{
		client: v.client,
		ID:     v.ID,
	}
}

// Get reads a certificate data from the specified file name.
func (a AutoCertCache) Get(ctx context.Context, name string) ([]byte, error) {
	// Validate parameters.
	if ctx == nil {
		return nil, errors.New("context is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if strings.Contains(name, "/") {
		return nil, errors.New("name cannot traverse directories")
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, a.ID) // nolint

	// Read secret.
	path := "secret/autocert/" + name
	vaultSecret, err := a.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}
	if vaultSecret == nil {
		if os.Getenv(env.Debug) != "" {
			logger.Verbose(ctx, "Secret does not exist: "+path)
		} else {
			logger.Verbose(ctx, "Secret does not exist.")
		}

		return nil, autocert.ErrCacheMiss
	}

	// Parse secret.
	secretString, ok := vaultSecret.Data["cert"].(string)
	if !ok {
		if os.Getenv(env.Debug) != "" {
			logger.Verbose(ctx, "Secret does not contain 'cert': "+path)
		} else {
			logger.Verbose(ctx, "Secret does not contain 'cert'.")
		}

		return nil, autocert.ErrCacheMiss
	}
	var secretData []byte
	secretData, err = base64.StdEncoding.DecodeString(secretString)
	if err != nil {
		logger.Fatal(ctx, "Unable to base64 decode secret: "+secretString)
	}

	// Log.
	if os.Getenv(env.Debug) != "" {
		logger.Verbose(ctx, "Read secret: "+path)
	} else {
		logger.Verbose(ctx, "Read secret.")
	}

	return secretData, nil
}

// Put writes the certificate data to the specified file name.
// The file will be created with 0600 permissions.
func (a AutoCertCache) Put(ctx context.Context, name string, data []byte) error {
	// Validate parameters.
	if ctx == nil {
		return errors.New("context is required")
	}
	if name == "" {
		return errors.New("name is required")
	}
	if strings.Contains(name, "/") {
		return errors.New("name cannot traverse directories")
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, a.ID) // nolint

	// Create secret.
	path := "secret/autocert/" + name
	secretData := map[string]interface{}{
		"cert": data,
	}
	_, err := a.client.Logical().Write(path, secretData)
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

// Delete removes the specified file name.
func (a AutoCertCache) Delete(ctx context.Context, name string) error {
	// Validate parameters.
	if ctx == nil {
		return errors.New("context is required")
	}
	if name == "" {
		return errors.New("name is required")
	}
	if strings.Contains(name, "/") {
		return errors.New("name cannot traverse directories")
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, a.ID) // nolint

	// Delete secret.
	path := "secret/autocert/" + name
	_, err := a.client.Logical().Delete(path)
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
