package auth

import (
	"net/http"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	customContext "github.com/IvanKondrashkov/go-shortener/internal/context"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
)

const (
	authCookie string = "auth"
)

func Authentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sc := securecookie.New(config.AuthKey, nil)
		var userID uuid.UUID

		if cookie, err := r.Cookie(authCookie); err == nil {
			if err = sc.Decode(authCookie, cookie.Value, &userID); err == nil {
				ctx := customContext.SetContextUserID(r.Context(), userID)
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		} else {
			userID, _ = uuid.NewRandom()
			if encoded, err := sc.Encode(authCookie, userID); err == nil {
				http.SetCookie(w, &http.Cookie{
					Name:  authCookie,
					Value: encoded,
				})
				ctx := customContext.SetContextUserID(r.Context(), userID)
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
	})
}
