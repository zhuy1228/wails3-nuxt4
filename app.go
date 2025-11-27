package main

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// App struct
type App struct {
	ctx context.Context
	app *application.App
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	a.ctx = ctx
	return nil
}
