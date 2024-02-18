package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/marcospereirampj/golang-open-telemetry/pkg/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/fx"
)

type (
	HTTPHandler interface {
		Method() string
		Pattern() string

		http.Handler
	}

	RouterParams struct {
		fx.In
		Handlers []HTTPHandler `group:"handlers"`
	}
)

func NewHTTPRouter(params RouterParams) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "healthy")
	})

	router.Route("/", func(r chi.Router) {
		for _, handler := range params.Handlers {
			r.Route(handler.Pattern(), func(subRoute chi.Router) {
				subRoute.Method(handler.Method(), "/", handler)
			})
		}
	})

	return router
}

func StartHTTPServer(lifecycle fx.Lifecycle, router *chi.Mux) {

	server := &http.Server{
		Addr:         ":8080",
		Handler:      otelhttp.NewHandler(router, "/"),
		ReadTimeout:  5000 * time.Second,
		WriteTimeout: 5000 * time.Second,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			errs := make(chan error, 2)

			otelShutdown, err := telemetry.SetupOTelSDK(context.Background())
			if err != nil {
				return err
			}

			defer func() {
				//nolint:typecheck
				err = errors.Join(err, otelShutdown(context.Background()))
			}()

			go func() {
				log.Printf("Listening and serving HTTP on %v", server.Addr)
				errs <- server.ListenAndServe()
			}()

			return err
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}

func Serve() {
	ServerDependencies := fx.Provide(
		NewHTTPRouter)

	app := fx.New(
		fx.Options(
			ServerDependencies,
			HomeModule,
		),
		fx.Invoke(StartHTTPServer),
	)

	app.Run()
}
