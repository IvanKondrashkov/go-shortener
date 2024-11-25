package main

import (
	"net/http"

	"github.com/IvanKondrashkov/go-shortener/internal/app"
	"github.com/IvanKondrashkov/go-shortener/storage"
)

func main() {
	memRepositoryImpl := storage.NewMemRepositoryImpl()
	app := app.NewApp(memRepositoryImpl)
	router := NewRouter(app)

	err := http.ListenAndServe(`:8080`, router)
	if err != nil {
		panic(err)
	}
}
