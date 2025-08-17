package auth

import (
	"net/http"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
)

// Authentication middleware проверяет/устанавливает аутентификацию пользователя
// Принимает:
// h - следующий обработчик в цепочке
// Возвращает обработчик с проверкой аутентификации
func Authentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sc := securecookie.New(config.AuthKey, nil)
		var userID uuid.UUID

		if cookie, err := r.Cookie(authCookie); err == nil {
			if err = sc.Decode(authCookie, cookie.Value, &userID); err == nil {
				ctx := SetContextUserID(r.Context(), userID)
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
				ctx := SetContextUserID(r.Context(), userID)
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
	})
}
