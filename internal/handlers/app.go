package handlers

type App struct {
	BaseURL    string
	repository Repository
}

func NewApp(BaseURL string, repository Repository) *App {
	return &App{
		BaseURL:    BaseURL,
		repository: repository,
	}
}
