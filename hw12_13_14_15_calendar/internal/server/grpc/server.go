package internalgrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/server"
	eventpb "github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/server/pb"
	"google.golang.org/grpc"
)

// Server модель сервера.
type Server struct {
	logger server.Logger
	config Config
	app    *app.App
}

// Config конфиг сервера.
type Config struct {
	Port int
}

// NewServer конструктор сервера.
func NewServer(logger server.Logger, app *app.App, config Config) *Server {
	return &Server{logger: logger, config: config, app: app}
}

// Start старт сервера.
func (s *Server) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		return fmt.Errorf("create grpc listener: %w", err)
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryServerRequestLoggerInterceptor(s.logger),
		),
	)

	eventpb.RegisterEventServiceServer(server, newEventService(s.app))
	s.logger.Log(fmt.Sprintf("starting grpc server on %s", lsn.Addr().String()))

	if err := server.Serve(lsn); err != nil {
		return fmt.Errorf("serve grpc connections: %w", err)
	}

	s.logger.Log("START SERVINT HTTP")

	return nil
}
