package transporthttp

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/WebCraftersGH/User-service/internal/middlewares"
	swaggerdocs "github.com/WebCraftersGH/User-service/internal/transport/http/docs"
	httphandlers "github.com/WebCraftersGH/User-service/internal/transport/http/handlers"
	"github.com/WebCraftersGH/User-service/pkg/logging"
)

func NewRouter(
	userHandler *httphandlers.UserHandler,
	healthHandler *httphandlers.HealthHandler,
	docsHandler *swaggerdocs.DocsHandler,
	authChecker middlewares.AuthChecker,
	logger logging.Logger,
	debugMode bool,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Use(middlewares.GenerateRequestID)
	router.Use(middlewares.LoggingMiddleware(logger))

	if debugMode {
		router.Use(middlewares.CORSMiddleware)

		router.HandleFunc("/swagger/openapi.json", docsHandler.ServeSpec).Methods(http.MethodGet)
		router.HandleFunc("/swagger/", docsHandler.ServeUI).Methods(http.MethodGet)
		router.HandleFunc("/swagger", docsHandler.RedirectToUI).Methods(http.MethodGet)
	}

	router.HandleFunc("/health", healthHandler.Check).Methods(http.MethodGet)

	protected := router.PathPrefix("/api/v1").Subrouter()
	protected.Use(middlewares.AuthFromToken(authChecker, logger))

	protected.HandleFunc("/users/me", userHandler.GetMe).Methods(http.MethodGet)
	protected.HandleFunc("/users/me", userHandler.DeleteUser).Methods(http.MethodDelete)
	protected.HandleFunc("/users/me", userHandler.UpdateUser).Methods(http.MethodPut)

	router.HandleFunc(
		"/api/v1/users/{uuid:[0-9a-fA-F-]+}",
		userHandler.GetUserByID,
	).Methods(http.MethodGet)

	return router
}
