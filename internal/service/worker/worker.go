package worker

import (
	"context"

	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	customContext "github.com/IvanKondrashkov/go-shortener/internal/service/middleware/auth"
	"go.uber.org/zap"
)

// SendDeleteBatchRequest отправляет задачу на пакетное удаление в очередь обработки
// Принимает:
// ctx - контекст для контроля времени выполнения
// event - событие удаления (пакет URL и ID пользователя)
func (w *Worker) SendDeleteBatchRequest(ctx context.Context, event models.DeleteEvent) {
	select {
	case w.resultCh <- event:
	case <-ctx.Done():
		return
	}
}

// RunJobDeleteBatch запускает воркер для обработки задач удаления
// Принимает:
// ctx - контекст для контроля времени выполнения
func (w *Worker) RunJobDeleteBatch(ctx context.Context) {
	defer w.wg.Done()
	for event := range w.resultCh {
		if event.UserID == nil {
			continue
		}
		ctx = customContext.SetContextUserID(ctx, *event.UserID)
		err := w.service.DeleteBatchByUserID(ctx, event.Batch)
		if err != nil && ctx.Err() == nil {
			w.errorCh <- err
		}
	}
}

// ErrorListener обрабатывает ошибки от воркеров
// Принимает:
// ctx - контекст для контроля времени выполнения
// zl - логгер для записи ошибок
// Возвращает канал для сигнализации завершения ErrorListener
func (w *Worker) ErrorListener(ctx context.Context, zl *logger.ZapLogger) <-chan struct{} {
	defer close(w.doneCh)

	for err := range w.errorCh {
		select {
		case <-ctx.Done():
			zl.Log.Debug("user delete batch error (shutdown)", zap.Error(err))
		default:
			zl.Log.Debug("user delete batch error", zap.Error(err))
		}
	}
	return w.doneCh
}

// Close останавливает воркеры и освобождает ресурсы
func (w *Worker) Close() {
	close(w.resultCh)
	w.wg.Wait()
	close(w.errorCh)

	<-w.doneCh
}
