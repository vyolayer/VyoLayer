package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(logger *logger.AppLogger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)
		st, _ := status.FromError(err)

		userIDStr := "anonymous"
		if userID, extractErr := ctxutil.ExtractIAMUserUUID(ctx); extractErr == nil && userID != uuid.Nil {
			userIDStr = userID.String()
		}

		logAttributes := []any{
			slog.String("method", info.FullMethod),
			slog.String("user_id", userIDStr),
			slog.Duration("latency", duration),
			slog.String("grpc_code", st.Code().String()),
		}

		if err != nil {
			logAttributes = append(logAttributes, slog.String("error", err.Error()))
			logger.Error("Error: ", logAttributes)
		} else {
			logger.Info("Request: ", logAttributes)
		}

		return resp, err
	}
}
