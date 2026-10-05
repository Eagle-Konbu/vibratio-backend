package main

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) error {
	slog.InfoContext(ctx, "generate-episode invoked")
	return nil
}

func main() {
	lambda.Start(handler)
}
