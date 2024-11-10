package accounts

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/michaelpeterswa/talvi/backend/internal/cockroach"
	"github.com/michaelpeterswa/talvi/backend/internal/dragonfly"
	"github.com/michaelpeterswa/talvi/backend/internal/encryption"
	"github.com/michaelpeterswa/talvi/backend/internal/util"
	"github.com/redis/go-redis/v9"
)

//go:embed queries/accounts/create_account.pgsql
var createAccountSQL string

//go:embed queries/accounts/get_account.pgsql
var getAccountSQL string

// //go:embed queries/accounts/get_accounts.pgsql
// var getAccountsSQL string

// //go:embed queries/accounts/update_account.pgsql
// var updateAccountSQL string

// //go:embed queries/accounts/delete_account.pgsql
// var deleteAccountSQL string

type AccountsClient struct {
	kv *dragonfly.DragonflyClient
	db *cockroach.CockroachClient

	aesClient *encryption.AESClient
}

type Account struct {
	ID                string    `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	Name              string    `json:"name"`
	Role              string    `json:"role"`
	Email             string    `json:"email"`
	Provider          string    `json:"provider"`
	EmailProviderHash string    `json:"email_provider_hash"`
}

func (a *Account) ToJSON() (string, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func NewAccountsClient(kv *dragonfly.DragonflyClient, db *cockroach.CockroachClient, aesClient *encryption.AESClient) *AccountsClient {
	return &AccountsClient{
		kv:        kv,
		db:        db,
		aesClient: aesClient,
	}
}

func (ac *AccountsClient) CreateAccount(ctx context.Context, name string, role string, email string, provider string) (bool, error) {
	_, err := ac.db.Client.Exec(ctx, createAccountSQL, name, role, email, provider, util.GenerateEmailProviderHash(email, provider))
	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		// is duplicate key on email_provider
		if pgErr.Code == "23505" {
			return true, nil
		}

		return false, pgErr
	}
	if err != nil {
		return false, err
	}

	account, err := ac.GetAccount(ctx, email, provider)
	if err != nil {
		return false, err
	}

	b, err := json.Marshal(account)
	if err != nil {
		return false, fmt.Errorf("error marshalling account to json: %w", err)
	}

	err = ac.kv.Client.Set(ctx, util.GenerateEmailProviderHash(email, provider), b, time.Hour*24).Err()
	if err != nil {
		return false, fmt.Errorf("error setting account in kv store: %w", err)
	}

	return false, nil
}

func (ac *AccountsClient) GetAccount(ctx context.Context, email string, provider string) (*Account, error) {
	accountJSON, err := ac.kv.Client.Get(ctx, util.GenerateEmailProviderHash(email, provider)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("error getting account from kv store: %w", err)
	} else if err == nil {
		var a Account
		err = json.NewDecoder(bytes.NewBufferString(accountJSON)).Decode(&a)
		if err != nil {
			return nil, fmt.Errorf("error decoding account from json: %w", err)
		}
		return &a, nil
	}

	row := ac.db.Client.QueryRow(ctx, getAccountSQL, util.GenerateEmailProviderHash(email, provider))
	var account Account
	err = row.Scan(&account.ID, &account.CreatedAt, &account.Name, &account.Role, &account.Email, &account.Provider, &account.EmailProviderHash)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	} else if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &account, nil
}
