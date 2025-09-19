package admin

import (
	"net/http"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/IvanKondrashkov/go-shortener/internal/utils/admin"
)

// TrustedSubnet middleware проверяет, принадлежит ли IP адрес указанной подсети
// Принимает:
// h - следующий обработчик в цепочке
// Возвращает обработчик с проверкой IP адреса
func TrustedSubnet(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.Header.Get("X-Real-IP")
		isTrusted, err := admin.IsIPInSubnet(clientIP, config.TrustedSubnet)
		if err != nil {
			http.Error(w, "is incorrect!", http.StatusInternalServerError)
			return
		}

		if !isTrusted {
			http.Error(w, "Access denied!", http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}
