package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/michaelpeterswa/talvi/backend/internal/accounts"
	"github.com/michaelpeterswa/talvi/backend/internal/middleware"
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

	jwt, err := getJWTFromRequestContext(r)
	if err != nil {
		ah.logger.Info("error getting jwt from request context", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if jwt.Email != account.Email || jwt.Provider != account.Provider {
		ah.logger.Info("error creating account: email and provider do not match jwt",
			zap.String("account_email", account.Email),
			zap.String("account_provider", account.Provider),
			zap.String("jwt_email", jwt.Email),
			zap.String("jwt_provider", jwt.Provider))
		w.WriteHeader(http.StatusUnauthorized)
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

	jwt, err := middleware.GetJWTFromRequestContext(r)
	if err != nil {
		ah.logger.Info("error getting jwt from request context", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if jwt.Email != email || jwt.Provider != provider {
		ah.logger.Info("error creating account: email and provider do not match jwt",
			zap.String("account_email", email),
			zap.String("account_provider", provider),
			zap.String("jwt_email", jwt.Email),
			zap.String("jwt_provider", jwt.Provider))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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

func (ah *AccountsHandler) Create2FA(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	provider := r.URL.Query().Get("provider")

	jwt, err := middleware.GetJWTFromRequestContext(r)
	if err != nil {
		ah.logger.Info("error getting jwt from request context", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if jwt.Email != email || jwt.Provider != provider {
		ah.logger.Info("error creating 2fa: email and provider do not match jwt",
			zap.String("2fa_email", email),
			zap.String("2fa_provider", provider),
			zap.String("jwt_email", jwt.Email),
			zap.String("jwt_provider", jwt.Provider))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	png, err := ah.accountsClient.Create2FA(r.Context(), email, provider)
	if err != nil {
		ah.logger.Info("error creating 2fa", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	outFile, err := os.Create("/etc/talvi/output/2fa.png")
	if err != nil {
		ah.logger.Info("error creating 2fa.png", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = outFile.Write(png)
	if err != nil {
		ah.logger.Info("error writing 2fa.png", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// func (ah *AccountsHandler) Get2FA(w http.ResponseWriter, r *http.Request) {
// 	email := r.URL.Query().Get("email")
// 	provider := r.URL.Query().Get("provider")

// 	jwt, err := middleware.GetJWTFromRequestContext(r)
// 	if err != nil {
// 		ah.logger.Info("error getting jwt from request context", zap.Error(err))
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	if jwt.Email != email || jwt.Provider != provider {
// 		ah.logger.Info("error getting 2fa: email and provider do not match jwt",
// 			zap.String("2fa_email", email),
// 			zap.String("2fa_provider", provider),
// 			zap.String("jwt_email", jwt.Email),
// 			zap.String("jwt_provider", jwt.Provider))
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	twoFactor, err := ah.accountsClient.Get2FA(r.Context(), email, provider)
// 	if err != nil {
// 		ah.logger.Info("error getting 2fa", zap.Error(err))
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	if twoFactor == nil {
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}

// 	res, err := json.Marshal(twoFactor)
// 	if err != nil {
// 		ah.logger.Info("error marshalling 2fa to json", zap.Error(err))
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	_, _ = w.Write([]byte(res))
// }

type Verify2FAResponse struct {
	Verified bool `json:"verified"`
}

func (ah *AccountsHandler) Verify2FA(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	provider := r.URL.Query().Get("provider")
	code := r.URL.Query().Get("code")

	jwt, err := middleware.GetJWTFromRequestContext(r)
	if err != nil {
		ah.logger.Info("error getting jwt from request context", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if jwt.Email != email || jwt.Provider != provider {
		ah.logger.Info("error verifying 2fa: email and provider do not match jwt",
			zap.String("2fa_email", email),
			zap.String("2fa_provider", provider),
			zap.String("jwt_email", jwt.Email),
			zap.String("jwt_provider", jwt.Provider))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	verified, err := ah.accountsClient.Verify2FA(r.Context(), email, provider, code)
	if err != nil {
		ah.logger.Info("error verifying 2fa", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
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
