package awssecretsmanager

import (
	"context"
	"errors"
)

// CreateToken creates a token.
func (a *AWSSecretsManager) CreateToken(ctx context.Context, id string, displayName string, numUses int, policies []string) (token string, err error) {
	return "", errors.New("not implemented")
}
