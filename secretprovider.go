// Package secretprovider provides methods for working with ISecretProvider.
package secretprovider

import (
	"context"
	"errors"
	"strings"

	"github.com/bertjohnson/secretprovider/awssecretsmanager"
	"github.com/bertjohnson/secretprovider/localfiles"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
	"github.com/bertjohnson/secretprovider/vault"
)

// Get returns a matching ISecretProvider for a secret store definition.
func Get(ctx context.Context, secretProvider *secretprovidertype.SecretProvider) (secretprovidertype.ISecretProvider, error) {
	// Validate parameters.
	if ctx == nil {
		return nil, errors.New("context is required")
	}
	if secretProvider == nil {
		return nil, errors.New("secret provider is required")
	}

	// Initialize based on the secret store type.
	switch strings.ToLower(secretProvider.Type) {
	case "awssecretsmanager":
		return awssecretsmanager.New(ctx, secretProvider)
	case "localfiles":
		return localfiles.New(ctx, secretProvider)
	case "vault":
		return vault.New(ctx, secretProvider)
	default:
		return nil, errors.New("unknown secret provider type: " + secretProvider.Type)
	}
}
