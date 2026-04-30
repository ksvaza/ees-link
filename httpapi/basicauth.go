package httpapi

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/ksvaza/ees-link/db"
	"github.com/ksvaza/ees-link/models"
	"github.com/pkg/errors"
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

/*
func BasicAuth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			errorHandler(w, errors.New("unauthorized"), http.StatusUnauthorized)
			return
		}

		u, err := lookupUser(username)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			errorHandler(w, errors.New("unauthorized"), http.StatusUnauthorized)
			return
		}

		if hashPassword(password, u.Salt) != u.PasswordHash {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			errorHandler(w, errors.New("unauthorized"), http.StatusUnauthorized)
			return
		}

		next(w, r, ps)
	}
}
*/

func Authenticate(r *http.Request) (*models.Account, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, errors.New("unauthorized")
	}

	database := db.RealDB{}
	a, err := database.GetAccountByUsername(r.Context(), username)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	if hashPassword(password, a.Salt) != a.Password {
		return nil, errors.New("unauthorized")
	}

	return &a, nil
}
