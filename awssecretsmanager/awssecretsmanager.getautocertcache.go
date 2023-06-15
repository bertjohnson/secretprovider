package awssecretsmanager

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	"github.com/bertjohnson/logger/types/env"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
	"golang.org/x/crypto/acme/autocert"
)

// AutoCertCache implements AutoCertCache using Vault.
type AutoCertCache struct {
	ID             string
	secretsManager *secretsmanager.SecretsManager
}

// GetAutoCertCache returns an autocert-compatible cache.
func (a *AWSSecretsManager) GetAutoCertCache(ctx context.Context) secretprovidertype.AutoCertCache {
	return AutoCertCache{
		ID:             a.ID,
		secretsManager: a.secretsManager,
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
	path := "autocert/" + name
	secret, err := a.secretsManager.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: &path,
	})
	if err != nil {
		return nil, err
	}
	if secret == nil {
		if os.Getenv(env.Debug) != "" {
			logger.Verbose(ctx, "Secret does not exist: "+path)
		} else {
			logger.Verbose(ctx, "Secret does not exist.")
		}

		return nil, autocert.ErrCacheMiss
	}

	// Parse secret.
	secretData := secret.SecretBinary

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
	path := "autocert/" + name
	_, err := a.secretsManager.PutSecretValue(&secretsmanager.PutSecretValueInput{
		SecretBinary: data,
		SecretId:     &path,
	})
	// If the secret does not exist, create it.
	if err != nil && strings.HasPrefix(err.Error(), "ResourceNotFoundException") {
		_, err = a.secretsManager.CreateSecret(&secretsmanager.CreateSecretInput{
			Name:         &path,
			SecretBinary: data,
		})
	}
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
	path := "autocert/" + name
	_, err := a.secretsManager.DeleteSecret(&secretsmanager.DeleteSecretInput{
		SecretId: &path,
	})
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
