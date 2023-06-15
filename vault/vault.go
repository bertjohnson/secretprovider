// Package vault hosts the Vault type.
package vault

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	"github.com/bertjohnson/logger/types/env"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
	vault "github.com/hashicorp/vault/api"
)

// Vault provides methods for interacting with Vault.
type Vault struct {
	ID string

	client      *vault.Client
	vaultConfig *vault.Config
}

// New creates a matching secret store implementation.
func New(ctx context.Context, secretProvider *secretprovidertype.SecretProvider) (*Vault, error) {
	// Validate parameters.
	if ctx == nil {
		return nil, errors.New("context is required")
	}
	if secretProvider == nil {
		return nil, errors.New("secret store configuration is required")
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, secretProvider.ID) // nolint

	// Log.
	logger.Verbose(ctx, "Creating Vault client.")

	// Inherit defaults.
	if secretProvider.ClientSecret == "" {
		secretProvider.ClientSecret = os.Getenv(env.SecretProviderClientSecret)
		if secretProvider.ClientSecret == "" {
			secretProvider.ClientSecret = os.Getenv(env.VaultToken)
			if secretProvider.ClientSecret == "" {
				return nil, errors.New("Vault token is required")
			}
		}
	}
	if len(secretProvider.UnsealShards) == 0 {
		unsealShards := os.Getenv(env.VaultUnsealShards)
		if unsealShards != "" {
			secretProvider.UnsealShards = strings.Split(unsealShards, ",")
			unsealShardsLength := len(secretProvider.UnsealShards)
			if unsealShardsLength > 0 {
				for i := 0; i < unsealShardsLength; i++ {
					secretProvider.UnsealShards[i] = strings.Trim(secretProvider.UnsealShards[i], " \t")
				}
			}
		}
	}
	if secretProvider.URI == "" {
		secretProvider.URI = os.Getenv(env.VaultAddr)
	}

	// Initialize Vault client.
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = secretProvider.URI

	// Configure TLS.
	skipVerify := strings.ToLower(os.Getenv(env.VaultSkipVerify))
	if boolVal, _ := strconv.ParseBool(skipVerify); boolVal { // nolint
		vaultTLSConfig := vault.TLSConfig{
			Insecure: true,
		}
		logger.Info(ctx, "Skipping Vault TLS certificate verification.")
		err := vaultConfig.ConfigureTLS(&vaultTLSConfig)
		if err != nil {
			return nil, err
		}
	}

	// Wait until the shard is unsealed.
	var err error
	vaultClient := Vault{
		ID:          secretProvider.ID,
		vaultConfig: vaultConfig,
	}
	logger.Info(ctx, "Unsealing Vault.")
	vaultClient.client, err = vault.NewClient(vaultClient.vaultConfig)
	if err != nil {
		return nil, err
	}
	var resp *vault.SealStatusResponse
	for _, unsealShard := range secretProvider.UnsealShards {
		if os.Getenv(env.Debug) != "" {
			logger.Verbose(ctx, "Unsealing using shard: "+unsealShard)
		}
		resp, err = vaultClient.client.Sys().Unseal(unsealShard)
		if err != nil {
			return nil, err
		}
	}
	if resp == nil {
		resp, err = vaultClient.client.Sys().SealStatus()
		if err != nil {
			return nil, err
		}
	}
	if resp.Sealed {
		return nil, errors.New("Vault is sealed")
	}

	logger.Info(ctx, "Unsealed Vault.")

	// Set the worker token.
	vaultClient.client.SetToken(secretProvider.ClientSecret)

	// Log.
	logger.Verbose(ctx, "Created Vault client.")

	return &vaultClient, nil
}
