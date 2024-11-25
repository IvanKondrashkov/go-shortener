package app

import (
	"github.com/IvanKondrashkov/go-shortener/storage"
)

type App struct {
	memRepository storage.MemRepository
}

func NewApp(memRepository *storage.MemRepositoryImpl) *App {
	return &App{
		memRepository: memRepository,
	}
}
