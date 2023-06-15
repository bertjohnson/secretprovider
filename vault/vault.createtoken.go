package vault

import (
	"context"
	"errors"

	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	vault "github.com/hashicorp/vault/api"
)

// CreateToken creates a token.
func (v *Vault) CreateToken(ctx context.Context, id string, displayName string, numUses int, policies []string) (token string, err error) {
	// Validate parameters.
	if ctx == nil {
		return "", errors.New("context is required")
	}
	if id == "" {
		return "", errors.New("token ID is required")
	}
	if displayName == "" {
		return "", errors.New("display name is required")
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, v.ID) // nolint

	// Create token.
	tokenCreateRequest := vault.TokenCreateRequest{
		DisplayName: displayName,
		ID:          id,
		NumUses:     numUses,
		Policies:    policies,
	}
	secret, err := v.client.Auth().Token().Create(&tokenCreateRequest)
	if err != nil {
		return "", err
	}

	// Log.
	logger.Info(ctx, "Created token.")

	return secret.Auth.ClientToken, nil
}
