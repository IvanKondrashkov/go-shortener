package handlers

type App struct {
	BaseURL    string
	repository repository
}

func NewApp(BaseURL string, repository repository) *App {
	return &App{
		BaseURL:    BaseURL,
		repository: repository,
	}
}
