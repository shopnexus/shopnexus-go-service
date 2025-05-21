package http

import (
	"fmt"
	"net/http"
	"shopnexus-go-service/internal/logger"
	"shopnexus-go-service/internal/server/http/response"
	"shopnexus-go-service/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init(mux *http.ServeMux, services *service.Services) error {
	handler := NewHandler(services)
	e := handler.SetupRoutes()

	// Wrap Echo with http.Handler
	mux.Handle("/", e)

	return nil
}

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) SetupRoutes() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Pre(middleware.AddTrailingSlash())

	// setup 404 handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			response.FromHTTPError(c.Response().Writer, he.Code)
		} else {
			response.FromHTTPError(c.Response().Writer, http.StatusInternalServerError)
		}
	}

	// Routes
	api := e.Group("/api")
	v1 := api.Group("/v1")
	vnpay := v1.Group("/vnpay-ipn")
	vnpay.GET("/", h.IPNVNPAY)
	vnpay.POST("/", h.IPNVNPAY)

	// Print routes for debugging
	logger.Log.Info("Registered routes:")
	for _, route := range e.Routes() {
		logger.Log.Info(fmt.Sprintf("  %s %s", route.Method, route.Path))
	}

	return e
}
