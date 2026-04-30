package httpapi

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"net/http"

	"github.com/ksvaza/ees-link/db"
	"github.com/ksvaza/ees-link/models"
	"github.com/sirupsen/logrus"
)

func hashPassword(password, salt string) string {
	h := sha512.New()
	h.Write([]byte(password + salt))
	return hex.EncodeToString(h.Sum(nil))
}

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
