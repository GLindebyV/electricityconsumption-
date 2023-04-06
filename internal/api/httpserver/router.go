package httpserver

import (
	"github.com/GLindebyV/electricityconsumption/internal/api/httpserver/consumptionupploader"
	"github.com/GLindebyV/electricityconsumption/internal/api/httpserver/health"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	router                         chi.Router
	healthController               *health.Controller
	consumptionupploaderController *consumptionupploader.Controller
}

func NewRouter(
	healthController *health.Controller,
	consumptionUpploaderController *consumptionupploader.Controller,
	server *Server,
) Router {
	return Router{
		router:                         server.GetRouter(),
		healthController:               healthController,
		consumptionupploaderController: consumptionUpploaderController,
	}
}

func (r Router) AddRoutes() {
	r.router.Route("/", func(rr chi.Router) {
		rr.Get("/health", r.healthController.HealthCheck)
		rr.Post("/uploadtest", r.consumptionupploaderController.UpploadTest)
		rr.Post("/consumption", r.consumptionupploaderController.UpploadConsumption)
	})
}
