package encryption_test

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/michaelpeterswa/talvi/backend/internal/encryption"
	"github.com/stretchr/testify/assert"
)

func TestEncryptAESGCMAndDecrypt(t *testing.T) {
	tests := []struct {
		Name      string
		Key       string
		Plaintext string
	}{
		{
			Name:      "Simple",
			Key:       "0936B5920B3E6FDDFEE77AE131C14385",
			Plaintext: "Hello, World!",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			aesConfig, err := encryption.NewAESConfig(tc.Key)
			assert.NoError(t, err)

			aesClient, err := encryption.NewAESClient(aesConfig)
			assert.NoError(t, err)

			ciphertextBytes, err := aesClient.Encrypt(ctx, []byte(tc.Plaintext))
			assert.NoError(t, err)

			plaintextBytes, err := aesClient.Decrypt(ctx, ciphertextBytes)
			assert.NoError(t, err)

			assert.Equal(t, tc.Plaintext, string(plaintextBytes))
		})
	}
}

func TestAESConfig(t *testing.T) {
	tests := []struct {
		Name string
		Key  string
	}{
		{
			Name: "Key Too Short",
			Key:  "0936B5920B3E6FDDFEE77",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := encryption.NewAESConfig(tc.Key)
			assert.ErrorIs(t, err, hex.ErrLength)
		})
	}
}
