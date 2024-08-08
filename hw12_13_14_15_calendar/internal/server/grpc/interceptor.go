package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// UnaryServerRequestLoggerInterceptor логирование входящих запросов.
func UnaryServerRequestLoggerInterceptor(logger server.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		t := time.Now()

		resp, err := handler(ctx, req)

		code := status.Code(err)

		_, offset := t.Zone()

		mData, _ := metadata.FromIncomingContext(ctx)

		p, _ := peer.FromContext(ctx)

		formatted := fmt.Sprintf("%02d/%s/%d:%02d:%02d:%02d +%04d",
			t.Day(), t.Month().String(), t.Year(),
			t.Hour(), t.Minute(), t.Second(), offset)

		str := fmt.Sprintf(`%s [%v] %s %s %v %v "%s"`,
			p.Addr.String(), formatted, info.FullMethod, "GRPC", code, time.Since(t).Milliseconds(), mData["user-agent"])

		logger.Log(str)

		return resp, err
	}
}
