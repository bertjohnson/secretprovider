package awssecretsmanager

import (
	"log"
	"testing"

	utilstrings "github.com/bertjohnson/util/strings"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestUpsertSecret tests UpsertSecret().
func TestUpsertSecret(t *testing.T) {
	secretID, err := utilstrings.HexToBase58(uuid.New().String())
	if err != nil {
		log.Fatalln(err.Error())
	}
	secretPath := "unittests/upsert/" + secretID
	writeSecret := map[string]interface{}{
		"a": 1,
		"b": "two",
		"c": []int{3},
	}
	err = awsSecretsManager.UpsertSecret(ctx, secretPath, writeSecret)
	assert.NoError(t, err)

	readSecret, err := awsSecretsManager.ReadSecret(ctx, secretPath)
	assert.NoError(t, err)
	if readSecret != nil && writeSecret != nil {
		assert.Equal(t, float64(1), readSecret.Data["a"])
		assert.Equal(t, writeSecret["b"], readSecret.Data["b"])
	}
}
