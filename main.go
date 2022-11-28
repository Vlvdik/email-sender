package main

import (
	"context"
	"mail-Sender/pkg/app"
)

func main() {
	sd, es := app.Init()
	ctx := context.Background()

	defer app.Close(ctx, &es)

	app.Start(&sd, &es)
}
