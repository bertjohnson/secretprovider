package vault

import (
	"testing"

	"encoding/json"
	"github.com/stretchr/testify/assert"
)

// TestReadSecret tests ReadSecret().
func TestReadSecret(t *testing.T) {
	secretPath := "secret/readsecret"
	writeSecret := map[string]interface{}{
		"a": 1,
		"b": "two",
		"c": []int{3},
	}
	err := vaultClient.UpsertSecret(ctx, secretPath, writeSecret)
	assert.NoError(t, err)

	// Read existing secret.
	readSecret, err := vaultClient.ReadSecret(ctx, secretPath)
	assert.NoError(t, err)
	if readSecret != nil && writeSecret != nil {
		readA, err := readSecret.Data["a"].(json.Number).Int64()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), readA)
		assert.Equal(t, writeSecret["b"], readSecret.Data["b"])
	}

	// Read nonexistent secret.
	readSecret, err = vaultClient.ReadSecret(ctx, "secret/nonexistent/path")
	assert.Error(t, err)
	assert.Equal(t, "not found", err)
}
