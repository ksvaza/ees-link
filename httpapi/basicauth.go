package httpapi

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/ksvaza/ees-link/db"
	"github.com/ksvaza/ees-link/models"
	"github.com/sirupsen/logrus"
)

func GenerateSalt(size int) (string, error) {
	if size <= 0 {
		size = 16
	}

	salt := make([]byte, size)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate random salt: %w", err)
	}

	return base64.RawStdEncoding.EncodeToString(salt), nil
}

type user struct {
	Username     string
	PasswordHash string
	Salt         string
}

func hashPassword(password, salt string) string {
	h := sha512.New()
	h.Write([]byte(password + salt))
	return hex.EncodeToString(h.Sum(nil))
}

// Context helpers
// ----------------------------------------------------------------

type contextKeyT string

var accountKey = contextKeyT("account")

func WithAccount(ctx context.Context, account *models.Account) context.Context {
	return context.WithValue(ctx, accountKey, account)
}

func GetAccount(ctx context.Context) *models.Account {
	account, ok := ctx.Value(accountKey).(*models.Account)
	if !ok {
		return nil
	}
	return account
}

// ----------------------------------------------------------------

func Authenticate(r *http.Request) *http.Request {
	ctx := r.Context()

	username, password, ok := r.BasicAuth()
	if !ok {
		return r
	}

	database := db.RealDB{}
	a, err := database.GetAccountByUsername(ctx, username)
	if err != nil {
		logrus.WithError(err).Error("Failed to get account by username")
		return r
	}

	if hashPassword(password, a.Salt) != a.Password {
		return r
	}

	ctx = WithAccount(ctx, &a)
	r = r.WithContext(ctx)

	return r
}
