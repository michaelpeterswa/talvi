package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/michaelpeterswa/talvi/backend/internal/accounts"
	"go.uber.org/zap"
)

type AccountsHandler struct {
	accountsClient *accounts.AccountsClient
	logger         *zap.Logger
}

func NewAccountsHandler(logger *zap.Logger, accountsClient *accounts.AccountsClient) *AccountsHandler {
	return &AccountsHandler{
		logger:         logger,
		accountsClient: accountsClient,
	}
}

type CreateAccountBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Provider string `json:"provider"`
}

func (ah *AccountsHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var account CreateAccountBody
	err = json.Unmarshal(body, &account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	exists, err := ah.accountsClient.CreateAccount(r.Context(), account.Name, account.Role, account.Email, account.Provider)
	if err != nil {
		ah.logger.Info("error creating account", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists {
		ah.logger.Info("account already exists")
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ah *AccountsHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	provider := r.URL.Query().Get("provider")

	account, err := ah.accountsClient.GetAccount(r.Context(), email, provider)
	if err != nil {
		ah.logger.Info("error getting account", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if account == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	accountJSON, err := account.ToJSON()
	if err != nil {
		ah.logger.Info("error marshalling account to json", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(accountJSON))
}

type Generate2FAResponse struct {
	Secret string `json:"secret"`
	Image  []byte `json:"image"`
}

func (ah *AccountsHandler) Generate2FA(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	provider := r.URL.Query().Get("provider")

	generated2FA, err := ah.accountsClient.Generate2FA(r.Context(), email, provider)
	if err != nil {
		ah.logger.Info("error creating 2fa", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json2FA, err := json.Marshal(Generate2FAResponse{
		Secret: generated2FA.Secret,
		Image:  generated2FA.Image,
	})
	if err != nil {
		ah.logger.Info("error marshalling 2fa to json", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(json2FA)
}

type Get2FAResponse struct {
	Status string `json:"status"`
}

func (ah *AccountsHandler) Get2FA(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	provider := r.URL.Query().Get("provider")

	twoFactor, err := ah.accountsClient.Get2FA(r.Context(), email, provider)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		ah.logger.Info("error getting 2fa", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if errors.Is(err, pgx.ErrNoRows) {
		json2FAStatus, err := json.Marshal(Get2FAResponse{Status: "not set up"})
		if err != nil {
			ah.logger.Info("error marshalling 2fa status to json", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(json2FAStatus)
		return
	}

	var resp Get2FAResponse
	if twoFactor.Enabled {
		resp.Status = "enabled"
	} else {
		resp.Status = "disabled"
	}

	res, err := json.Marshal(resp)
	if err != nil {
		ah.logger.Info("error marshalling 2fa resp to json", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(res))
}

type Create2FABody struct {
	Secret string
}

func (ah *AccountsHandler) Create2FA(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	provider := r.URL.Query().Get("provider")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var create2FABody Create2FABody
	err = json.Unmarshal(body, &create2FABody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = ah.accountsClient.Create2FA(r.Context(), email, provider, create2FABody.Secret)
	if err != nil {
		ah.logger.Info("error creating 2fa", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ah *AccountsHandler) Delete2FA(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	provider := r.URL.Query().Get("provider")

	err := ah.accountsClient.Delete2FA(r.Context(), email, provider)
	if err != nil {
		ah.logger.Info("error deleting 2fa", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type Verify2FAResponse struct {
	Verified bool `json:"verified"`
}

func (ah *AccountsHandler) Verify2FA(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	provider := r.URL.Query().Get("provider")
	code := r.URL.Query().Get("code")

	verified, err := ah.accountsClient.Verify2FA(r.Context(), email, provider, code)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		ah.logger.Info("error verifying 2fa", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if errors.Is(err, pgx.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := json.Marshal(Verify2FAResponse{Verified: verified})
	if err != nil {
		ah.logger.Info("error marshalling verify 2fa response to json", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(res)
}

func (ah *AccountsHandler) Validate2FA(w http.ResponseWriter, r *http.Request) {
	secret := r.URL.Query().Get("secret")
	code := r.URL.Query().Get("code")

	valid, err := ah.accountsClient.Validate2FASecretCode(r.Context(), secret, code)
	if err != nil {
		ah.logger.Info("error validating 2fa", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(Verify2FAResponse{Verified: valid})
	if err != nil {
		ah.logger.Info("error marshalling validate 2fa response to json", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(res)
}
