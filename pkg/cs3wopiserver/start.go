package cs3wopiserver

import (
	"context"
	"time"

	"github.com/wkloucek/cs3-wopi-server/pkg/internal/app"
)

func Start() error {
	ctx := context.Background()

	app, err := app.New()
	if err != nil {
		return err
	}

	if err := app.RegisterOcisService(ctx); err != nil {
		return err
	}

	if err := app.WopiDiscovery(ctx); err != nil {
		return err
	}

	if err := app.GetCS3apiClient(); err != nil {
		return err
	}

	if err := app.RegisterDemoApp(ctx); err != nil {
		return err
	}

	if err := app.GRPCServer(ctx); err != nil {
		return err
	}

	if err := app.HTTPServer(ctx); err != nil {
		return err
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
