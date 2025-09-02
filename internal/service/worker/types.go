package worker

import (
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
}

// NewWorker создает новый пул воркеров для обработки удаления URL
// Принимает:
// - workerCount: количество воркеров
// - zl: логгер
// - s: сервис для операций с URL
// Возвращает инициализированный Worker
func NewWorker(workerCount int, zl *logger.ZapLogger, s *service.Service) *Worker {
	w := &Worker{
		service:  s,
		resultCh: make(chan models.DeleteEvent, bufCh),
		errorCh:  make(chan error, bufCh),
	}

	go w.ErrorListener(zl)

	for i := 0; i < workerCount; i++ {
		w.wg.Add(1)
		go w.RunJobDeleteBatch()
	}
	return w
}
