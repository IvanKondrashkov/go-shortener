package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	customContext "github.com/IvanKondrashkov/go-shortener/internal/service/middleware/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetStats(t *testing.T) {
	tc := NewSuite(t)
	tests := []struct {
		name          string
		trustedSubnet string
		xRealIP       string
		status        int
		want          []byte
	}{
		{
			name:          "ip no trusted subnet",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "10.0.0.1",
			status:        http.StatusForbidden,
			want:          []byte("Access denied!\n"),
		},
		{
			name:          "ok",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "192.168.1.100",
			status:        http.StatusOK,
			want:          []byte("{\"urls\":0,\"users\":0}\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.app.URL+"api/internal/stats", nil)
			req.Header.Set("X-Real-IP", tt.xRealIP)

			handler := NewHandler(nil, tc.app)
			router := NewRouter(handler)

			ctx := customContext.SetContextUserID(req.Context(), uuid.New())
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.status, w.Code)
			assert.Equal(t, tt.want, w.Body.Bytes())
		})
	}
}
