package server

import (
	"github.com/marcospereirampj/golang-open-telemetry/otlp/handlers/get/home"
	"go.uber.org/fx"
)

type HomeHandlerOutput struct {
	fx.Out
	Handler HTTPHandler `group:"handlers"`
}

func NewHomeServiceHandler() HomeHandlerOutput {
	return HomeHandlerOutput{
		Handler: home.NewHomeServiceHandler(),
	}
}

var HomeModule = fx.Provide(
	NewHomeServiceHandler,
)
