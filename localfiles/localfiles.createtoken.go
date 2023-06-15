package localfiles

import (
	"context"

	"github.com/bertjohnson/logger"
)

// CreateToken creates a token.
func (l *LocalFiles) CreateToken(ctx context.Context, id string, displayName string, numUses int, policies []string) (token string, err error) {
	// Log.
	logger.Info(ctx, "Created token.")

	return "", nil
}
