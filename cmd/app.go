// Package cmd implements shotty's CLI command handlers.
package cmd

import (
	"github.com/vdemeester/shotty/internal/config"
	"github.com/vdemeester/shotty/internal/ext"
)

// App holds shared dependencies for all commands.
type App struct {
	Tools  *ext.Tools
	Config *config.Config
}

// NewApp creates an App with default (real) dependencies.
func NewApp() *App {
	return &App{
		Tools:  ext.DefaultTools(),
		Config: config.New(),
	}
}
