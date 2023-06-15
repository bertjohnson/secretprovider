package localfiles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestReadWorkerSecret tests WriteSecret() and ReadSecret().
func TestReadSecret(t *testing.T) {
	secretPath := "secret/workersecret"
	writeSecret := map[string]interface{}{
		"a": float64(1),
		"b": "two",
		"c": []int{3},
	}
	err := localFilesClient.UpsertSecret(ctx, secretPath, writeSecret)
	assert.NoError(t, err)

	readSecret, err := localFilesClient.ReadSecret(ctx, secretPath)
	assert.NoError(t, err)
	if readSecret != nil && writeSecret != nil {
		readA, ok := readSecret.Data["a"]
		assert.True(t, ok)
		readAInt, ok := readA.(float64)
		assert.True(t, ok)
		assert.Equal(t, readAInt, writeSecret["a"])
	}
}
