package vault

import (
	"testing"

	"encoding/json"
	"github.com/stretchr/testify/assert"
)

// TestUpsertSecret tests UpsertSecret().
func TestUpsertSecret(t *testing.T) {
	secretPath := "secret/upsertsecret"
	writeSecret := map[string]interface{}{
		"a": 1,
		"b": "two",
		"c": []int{3},
	}
	err := vaultClient.UpsertSecret(ctx, secretPath, writeSecret)
	assert.NoError(t, err)

	readSecret, err := vaultClient.ReadSecret(ctx, secretPath)
	assert.NoError(t, err)
	if readSecret != nil && writeSecret != nil {
		readA, err := readSecret.Data["a"].(json.Number).Int64()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), readA)
		assert.Equal(t, writeSecret["b"], readSecret.Data["b"])
	}
}
