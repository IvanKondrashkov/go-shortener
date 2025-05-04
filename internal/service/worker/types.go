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

type Worker struct {
	wg       sync.WaitGroup
	service  *service.Service
	resultCh chan models.DeleteEvent
	errorCh  chan error
}

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
