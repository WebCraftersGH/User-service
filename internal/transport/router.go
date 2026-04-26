package transporthttp

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/WebCraftersGH/User-service/internal/middlewares"
	swaggerdocs "github.com/WebCraftersGH/User-service/internal/transport/http/docs"
	httphandlers "github.com/WebCraftersGH/User-service/internal/transport/http/handlers"
)

func NewRouter(
	userHandler *httphandlers.UserHandler,
	healthHandler *httphandlers.HealthHandler,
	docsHandler *swaggerdocs.DocsHandler,
	authChecker middlewares.AuthChecker,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// DocsHandlers
	router.HandleFunc("/swagger/openapi.json", docsHandler.ServeSpec).Methods(http.MethodGet)
	router.HandleFunc("/swagger/", docsHandler.ServeUI).Methods(http.MethodGet)
	router.HandleFunc("/swagger", docsHandler.RedirectToUI).Methods(http.MethodGet)

	//Health
	router.HandleFunc("/health", healthHandler.Check).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/users/{uuid}", userHandler.GetUserByID).Methods(http.MethodGet)

	protected := router.PathPrefix("/api/v1").Subrouter()

	//auth-middleware
	protected.Use(middlewares.AuthFromToken(authChecker))

	protected.HandleFunc("/users/me", userHandler.GetMe).Methods(http.MethodGet)
	protected.HandleFunc("/users/me", userHandler.DeleteUser).Methods(http.MethodDelete)
	protected.HandleFunc("/users/me", userHandler.UpdateUser).Methods(http.MethodPut)

	return router
}
