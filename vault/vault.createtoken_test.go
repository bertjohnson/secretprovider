package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateToken tests CreateToken.
func TestCreateToken(t *testing.T) {
	// Create token for account.
	accountID := "5c47ad37-7052-4e65-a0db-3ef2ac235e22"
	displayName := "Test User"
	numUses := 2
	policies := []string{
		"path \"secret/*\" { capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"] }",
	}
	accountToken, err := vaultClient.CreateToken(ctx, accountID, displayName, numUses, policies)
	assert.NotEqual(t, "", accountToken)
	assert.NoError(t, err)
}
