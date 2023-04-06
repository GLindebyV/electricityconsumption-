package app

import (
	"context"
	"net/http"

	"github.com/GLindebyV/electricityconsumption/internal/api/httpserver"
	"github.com/GLindebyV/electricityconsumption/internal/api/httpserver/consumptionupploader"
	"github.com/GLindebyV/electricityconsumption/internal/api/httpserver/health"
	"github.com/GLindebyV/electricityconsumption/internal/gateway/elpriserjustnu"
	"github.com/GLindebyV/electricityconsumption/internal/service"
	"go.uber.org/fx"
)

type App struct {
	*fx.App
}

func New(opts ...fx.Option) App {
	return App{
		App: fx.New(
			fx.Provide(
				httpserver.New,
				health.NewController,
				consumptionupploader.NewController,
				httpserver.NewRouter,
				NewHttpClient,
				fx.Annotate(
					elpriserjustnu.NewClient,
					fx.As(new(service.ElectricityPrice)),
				),
				fx.Annotate(
					service.NewPricingService,
					fx.As(new(consumptionupploader.PriceService)),
				),
			),
			fx.Options(opts...),
			fx.Invoke(func(router httpserver.Router) {
				router.AddRoutes()
			}),
			fx.Invoke(startHttpServer),
		)}
}

func startHttpServer(lc fx.Lifecycle, server *httpserver.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.Start(":8080"); err != nil {
					panic(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}

func NewHttpClient() *http.Client {
	return &http.Client{}
}
