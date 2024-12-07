package app

import (
	"github.com/IvanKondrashkov/go-shortener/storage"
)

type App struct {
	BaseURL       string
	memRepository storage.MemRepository
}

func NewApp(BaseURL string, memRepository *storage.MemRepositoryImpl) *App {
	return &App{
		BaseURL:       BaseURL,
		memRepository: memRepository,
	}
}
