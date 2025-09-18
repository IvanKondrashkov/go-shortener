package worker

import (
	"context"
	"sync"

	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
)

const (
	bufCh = 100
)

// Worker - структура для фоновой обработки задач удаления URL
type Worker struct {
	wg       sync.WaitGroup          // Группа ожидания завершения воркеров
	service  *service.Service        // Сервис для операций с URL
	resultCh chan models.DeleteEvent // Канал для задач удаления
	errorCh  chan error              // Канал для ошибок
	doneCh   chan struct{}           // Канал для сигнализации завершения ErrorListener
}

// NewWorker создает новый пул воркеров для обработки удаления URL
// Принимает:
// - ctx: контекст для контроля времени выполнения
// - workerCount: количество воркеров
// - zl: логгер
// - s: сервис для операций с URL
// Возвращает инициализированный Worker
func NewWorker(ctx context.Context, workerCount int, zl *logger.ZapLogger, s *service.Service) *Worker {
	w := &Worker{
		service:  s,
		resultCh: make(chan models.DeleteEvent, bufCh),
		errorCh:  make(chan error, bufCh),
		doneCh:   make(chan struct{}),
	}

	go w.ErrorListener(ctx, zl)

	for i := 0; i < workerCount; i++ {
		w.wg.Add(1)
		go w.RunJobDeleteBatch(ctx)
	}
	return w
}
