package main

import (
	"context"

	"github.com/multimoml/qr-generator/internal/server"
)

func main() {
	ctx := context.Background()
	server.Run(ctx)
}
