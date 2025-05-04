package worker

import (
	"context"

	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	customContext "github.com/IvanKondrashkov/go-shortener/internal/service/middleware/auth"
	"go.uber.org/zap"
)

func (w *Worker) SendDeleteBatchRequest(ctx context.Context, event models.DeleteEvent) {
	select {
	case w.resultCh <- event:
	case <-ctx.Done():
		return
	}
}

func (w *Worker) RunJobDeleteBatch() {
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

func (w *Worker) ErrorListener(zl *logger.ZapLogger) {
	for err := range w.errorCh {
		zl.Log.Debug("user delete batch error", zap.Error(err))
	}
}

func (w *Worker) Close() {
	close(w.resultCh)
	w.wg.Wait()
	close(w.errorCh)
}
