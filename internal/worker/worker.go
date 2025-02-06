package worker

import (
	"context"
	"sync"

	customContext "github.com/IvanKondrashkov/go-shortener/internal/context"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	"github.com/IvanKondrashkov/go-shortener/internal/service"

	"go.uber.org/zap"
)

const (
	bufCh = 100
)

type Worker struct {
	wg       sync.WaitGroup
	service  *service.Service
	resultCh chan models.DeleteEvent
	errorCh  chan error
}

func NewWorker(workerCount int, zl *logger.ZapLogger, newService *service.Service) *Worker {
	w := &Worker{
		service:  newService,
		resultCh: make(chan models.DeleteEvent, bufCh),
		errorCh:  make(chan error, bufCh),
	}

	for i := 0; i < workerCount; i++ {
		w.wg.Add(1)
		go w.errorListener(zl)
		go w.runJobDeleteBatch()
	}
	return w
}

func (w *Worker) SendDeleteBatchRequest(ctx context.Context, event models.DeleteEvent) {
	select {
	case w.resultCh <- event:
	case <-ctx.Done():
		return
	}
}

func (w *Worker) Close() {
	close(w.resultCh)
	w.wg.Wait()
	close(w.errorCh)
}

func (w *Worker) runJobDeleteBatch() {
	defer w.wg.Done()
	for event := range w.resultCh {
		if event.UserID == nil {
			continue
		}
		ctx := customContext.SetContextUserID(context.Background(), *event.UserID)
		err := w.service.DeleteBatchByUserID(ctx, event.Batch)
		if err != nil {
			w.errorCh <- err
		}
	}
}

func (w *Worker) errorListener(zl *logger.ZapLogger) {
	for err := range w.errorCh {
		zl.Log.Debug("user delete batch error", zap.Error(err))
	}
}
