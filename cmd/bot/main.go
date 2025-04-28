package main

import (
	"context"
	"log"
	"weatherbot/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.New(ctx)
	if err != nil {
		log.Fatalf("failed init app: %s", err.Error())
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("failed run app: %s", err.Error())
	}
}
