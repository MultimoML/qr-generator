package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/multimoml/qr-generator/internal/server"
)

// @title QR Generator API
// @version 1.0.0
// @host localhost:6002
// @BasePath /qr
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Gracefully exit on SIGINT and SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go server.Run(ctx)

	<-sigChan
	cancel()
}
