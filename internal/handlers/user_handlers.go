package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/IvanKondrashkov/go-shortener/internal/models"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
	customContext "github.com/IvanKondrashkov/go-shortener/internal/service/middleware/auth"
	"github.com/google/uuid"
)

func (app *App) GetAllURLByUserID(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	respDto, err := app.service.GetAllByUserID(req.Context())
	if err != nil && errors.Is(err, service.ErrUserUnauthorized) {
		res.WriteHeader(http.StatusUnauthorized)
		_, _ = res.Write([]byte("User unauthorized!"))
		return
	}

	if len(respDto) == 0 {
		res.WriteHeader(http.StatusNoContent)
		_, _ = res.Write([]byte("Urls by user id not found!"))
		return
	}

	res.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(res)
	err = enc.Encode(&respDto)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Response is invalidate!"))
		return
	}
}

func (app *App) DeleteBatchByUserID(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var reqDto []uuid.UUID
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&reqDto)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Body is invalidate!"))
		return
	}
	defer req.Body.Close()

	event := models.DeleteEvent{
		Batch:  reqDto,
		UserID: customContext.GetContextUserID(req.Context()),
	}

	go app.worker.SendDeleteBatchRequest(context.Background(), event)
	res.WriteHeader(http.StatusAccepted)
}
