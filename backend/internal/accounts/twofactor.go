package accounts

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/hex"
	"fmt"
	"image/png"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pquerna/otp/totp"
	"go.opentelemetry.io/otel"

	"github.com/michaelpeterswa/talvi/backend/internal/util"
)

var (
	tracer = otel.Tracer("github.com/michaelpeterswa/talvi/backend/internal/accounts")
)

//go:embed queries/twofactor/create_twofactor.pgsql
var createTwofactorSQL string

//go:embed queries/twofactor/get_twofactor.pgsql
var getTwofactorSQL string

//go:embed queries/twofactor/delete_twofactor.pgsql
var deleteTwofactorSQL string

//go:embed queries/twofactor/update_twofactor.pgsql
var updateTwofactorSQL string

type TwoFactor struct {
	ID                string    `json:"id"`
	ParentAccountHash string    `json:"parent_account_hash"`
	CreatedAt         time.Time `json:"created_at"`
	Secret            string    `json:"secret"`
	Enabled           bool      `json:"enabled"`
}

type Generated2FA struct {
	Secret string
	Image  []byte
}

func (ac *AccountsClient) Generate2FA(ctx context.Context, email string, provider string) (*Generated2FA, error) {
	_, span := tracer.Start(ctx, "Generate2FA")
	defer span.End()

	opts := totp.GenerateOpts{
		Issuer:      "talvi",
		AccountName: email,
	}

	key, err := totp.Generate(opts)
	if err != nil {
		return nil, fmt.Errorf("error generating totp key: %w", err)
	}

	var pngBuffer bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, fmt.Errorf("error generating totp key image: %w", err)
	}
	err = png.Encode(&pngBuffer, img)
	if err != nil {
		return nil, fmt.Errorf("error encoding totp key to png: %w", err)
	}

	return &Generated2FA{
		Secret: key.Secret(),
		Image:  pngBuffer.Bytes(),
	}, nil
}

func (ac *AccountsClient) Create2FA(ctx context.Context, email string, provider string, secret string) error {
	traceCtx, span := tracer.Start(ctx, "Create2FA")
	defer span.End()

	ciphertext, err := ac.aesClient.Encrypt(traceCtx, []byte(secret))
	if err != nil {
		return fmt.Errorf("error encrypting totp key: %w", err)
	}
	hexCiphertext := hex.EncodeToString(ciphertext)

	_, err = ac.db.Client.Exec(ctx, createTwofactorSQL, util.GenerateEmailProviderHash(email, provider), hexCiphertext, true)
	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		// is duplicate key on email_provider
		if pgErr.Code == "23505" {
			return fmt.Errorf("error creating 2fa: %w", err)
		}

		return pgErr
	}
	if err != nil {
		return fmt.Errorf("error exec create query: %w", err)
	}

	return nil
}

func (ac *AccountsClient) Get2FA(ctx context.Context, email string, provider string) (*TwoFactor, error) {
	var tf TwoFactor
	row := ac.db.Client.QueryRow(ctx, getTwofactorSQL, util.GenerateEmailProviderHash(email, provider))
	err := row.Scan(&tf.ID, &tf.ParentAccountHash, &tf.CreatedAt, &tf.Secret, &tf.Enabled)
	if err != nil {
		return nil, fmt.Errorf("error scanning row: %w", err)
	}

	return &tf, nil
}

func (ac *AccountsClient) Update2FA(ctx context.Context, email string, provider string, enabled bool) error {
	_, err := ac.db.Client.Exec(ctx, updateTwofactorSQL, util.GenerateEmailProviderHash(email, provider), enabled)
	if err != nil {
		return fmt.Errorf("error exec update query: %w", err)
	}

	return nil
}

func (ac *AccountsClient) Delete2FA(ctx context.Context, email string, provider string) error {
	fmt.Println("deleteTwofactorSQL", deleteTwofactorSQL)

	_, err := ac.db.Client.Exec(ctx, deleteTwofactorSQL, util.GenerateEmailProviderHash(email, provider))
	if err != nil {
		return fmt.Errorf("error exec delete query: %w", err)
	}

	return nil
}

func (ac *AccountsClient) Validate2FASecretCode(ctx context.Context, secret string, code string) (bool, error) {
	valid := totp.Validate(code, secret)
	if !valid {
		return false, nil
	}

	return true, nil
}

func (ac *AccountsClient) Verify2FA(ctx context.Context, email string, provider string, code string) (bool, error) {
	tf, err := ac.Get2FA(ctx, email, provider)
	if err != nil {
		return false, fmt.Errorf("error getting 2fa: %w", err)
	}

	decodedCiphertext, err := hex.DecodeString(tf.Secret)
	if err != nil {
		return false, fmt.Errorf("error hex decoding secret: %w", err)
	}

	plaintext, err := ac.aesClient.Decrypt(ctx, decodedCiphertext)
	if err != nil {
		return false, fmt.Errorf("error decrypting secret: %w", err)
	}

	valid := totp.Validate(code, string(plaintext))
	if !valid {
		return false, nil
	}

	return true, nil
}
