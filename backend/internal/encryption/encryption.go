package encryption

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"go.opentelemetry.io/otel"
)

var (
	tracer = otel.Tracer("github.com/michaelpeterswa/talvi/backend/internal/encryption")
)

type AESConfig struct {
	Key []byte `json:"key"`
}

func NewAESConfig(key string) (*AESConfig, error) {
	decodedKey, err := hex.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("error decoding key: %w", err)
	}
	return &AESConfig{
		Key: decodedKey,
	}, nil
}

type AESClient struct {
	gcm cipher.AEAD
}

func NewAESClient(ac *AESConfig) (*AESClient, error) {
	block, err := aes.NewCipher(ac.Key)
	if err != nil {
		return nil, fmt.Errorf("error creating new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating new gcm: %w", err)
	}

	return &AESClient{
		gcm: gcm,
	}, nil
}

func (ac *AESClient) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	_, span := tracer.Start(ctx, "Encrypt")
	defer span.End()

	nonce := make([]byte, ac.gcm.NonceSize())
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("error reading random bytes: %w", err)
	}
	return ac.gcm.Seal(nonce, nonce, data, nil), nil
}

func (ac *AESClient) Decrypt(ctx context.Context, data []byte) ([]byte, error) {
	_, span := tracer.Start(ctx, "Decrypt")
	defer span.End()

	nonce := data[:ac.gcm.NonceSize()]
	data = data[ac.gcm.NonceSize():]
	plaintext, err := ac.gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %w", err)
	}

	return plaintext, nil
}
