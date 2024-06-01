package app

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface { // TODO
}

type Storage interface { // TODO
}

func New(logger Logger, storage Storage) *App {
	return &App{Logger: logger, Storage: storage}
}
