package types

import (
	"context"
)

// ISecretProvider contains methods used to interface with secrets.
type ISecretProvider interface {
	GetAutoCertCache(ctx context.Context) AutoCertCache
	CreateToken(ctx context.Context, id string, displayName string, numUses int, policies []string) (token string, err error)
	DeleteSecret(ctx context.Context, path string) error
	ListSecrets(ctx context.Context, pathChannel chan string, errorChannel chan error)
	ReadAllSecrets(ctx context.Context, secretChannel chan *Secret, errorChannel chan error)
	ReadSecret(ctx context.Context, path string) (secret *Secret, err error)
	UpsertSecret(ctx context.Context, path string, data map[string]interface{}) error
}
