package httpapi

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
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

func lookupUser(username string) (*user, error) {
	// TODO: replace with actual DB lookup
	// return db.GetUserByUsername(ctx, username)

	stubSalt := "randomsalt123"
	stubHash := hashPassword("password123", stubSalt)

	if username == "admin" {
		return &user{
			Username:     "admin",
			PasswordHash: stubHash,
			Salt:         stubSalt,
		}, nil
	}

	return nil, errors.New("user not found")
}

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
