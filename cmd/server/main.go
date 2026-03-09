package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"gorankd/internal/cache"
	"gorankd/internal/ranking"
	"gorankd/internal/server"
	"gorankd/internal/store"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	c := cache.NewRedisCache()
	s := store.NewSpannerStore()
	rankingSvc := ranking.NewService(c, s)
	_ = server.NewGRPCServer(rankingSvc)

	grpcServer := grpc.NewServer()
	// TODO: register ranking service with grpcServer after proto-gen

	addr := ":50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("failed to listen", "addr", addr, "error", err)
		os.Exit(1)
	}

	go func() {
		slog.Info("gRPC server starting", "addr", addr)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server failed", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down gRPC server")
	grpcServer.GracefulStop()
	slog.Info("server stopped")
}
