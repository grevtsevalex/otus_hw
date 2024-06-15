package server

// Logger тип логгера.
type Logger interface {
	Log(msg string)
	Error(msg string)
}

type Server interface {
	Start()
	Stop()
}
