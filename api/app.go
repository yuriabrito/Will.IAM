package api

import (
	"fmt"
	"net"
	"net/http"

	"github.com/ghostec/Will.IAM/constants"
	"github.com/ghostec/Will.IAM/oauth2"
	"github.com/ghostec/Will.IAM/repositories"
	"github.com/ghostec/Will.IAM/usecases"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/topfreegames/extensions/middleware"
)

// App struct
type App struct {
	address         string
	config          *viper.Viper
	logger          logrus.FieldLogger
	router          *mux.Router
	server          *http.Server
	metricsReporter middleware.MetricsReporter
	storage         *repositories.Storage
	oauth2Provider  oauth2.Provider
}

// NewApp creates a new app
func NewApp(
	host string, port int, config *viper.Viper, logger logrus.FieldLogger,
	storageOrNil *repositories.Storage,
) (*App, error) {
	mr, err := middleware.NewDogStatsD(config)
	if err != nil {
		return nil, err
	}
	if storageOrNil == nil {
		storageOrNil = repositories.NewStorage()
	}
	a := &App{
		config:          config,
		address:         fmt.Sprintf("%s:%d", host, port),
		logger:          logger,
		metricsReporter: mr,
		storage:         storageOrNil,
	}
	err = a.configureApp()
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) configureApp() error {
	err := a.configurePG()
	if err != nil {
		return err
	}

	a.configureGoogleOAuth2Provider()
	a.configureServer()

	return nil
}

func (a *App) configureServer() {
	a.router = a.GetRouter()
	a.server = &http.Server{Addr: a.address, Handler: wrapHandlerWithResponseWriter(a.router)}
}

func (a *App) configurePG() error {
	if a.storage != nil && a.storage.PG != nil {
		return nil
	}
	return a.storage.ConfigurePG(a.config)
}

func (a *App) configureGoogleOAuth2Provider() {
	tokensRepo := repositories.NewTokens(a.storage)
	google := oauth2.NewGoogle(oauth2.GoogleConfig{
		ClientID:      a.config.GetString("oauth2.google.clientId"),
		ClientSecret:  a.config.GetString("oauth2.google.clientSecret"),
		RedirectURL:   a.config.GetString("oauth2.google.redirectUrl"),
		HostedDomains: a.config.GetStringSlice("oauth2.google.hostedDomains"),
	}, tokensRepo)
	a.oauth2Provider = google
}

// SetOAuth2Provider sets a provider in App
func (a *App) SetOAuth2Provider(provider oauth2.Provider) {
	a.oauth2Provider = provider
}

// GetRouter returns App's *mux.Router reference
func (a *App) GetRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.Version(constants.AppInfo.Version))
	r.Use(middleware.Logging(a.logger))
	r.Use(middleware.Metrics(a.metricsReporter))

	r.HandleFunc("/healthcheck", healthcheckHandler(
		repositories.NewHealthcheck(a.storage),
	)).Methods("GET").Name("healthcheck")

	serviceAccountsRepo := repositories.NewServiceAccounts(a.storage)
	rolesRepo := repositories.NewRoles(a.storage)
	permissionsRepo := repositories.NewPermissions(a.storage)
	serviceAccountsUseCase := usecases.NewServiceAccounts(
		serviceAccountsRepo, rolesRepo, permissionsRepo, a.oauth2Provider,
	)

	r.HandleFunc(
		"/service_accounts/{id}/permissions",
		serviceAccountsHasPermissionHandler(serviceAccountsUseCase),
	).
		Methods("GET").Name("serviceAccountsHasPermission")

	r.HandleFunc(
		"/authentication/build_url",
		authenticationBuildURLHandler(a.oauth2Provider),
	).
		Methods("GET").Name("authenticationBuildURLHandler")

	r.HandleFunc(
		"/authentication",
		authenticationHandler(a.oauth2Provider),
	).
		Methods("GET").Name("authenticationBuildURLHandler")

	r.HandleFunc(
		"/authentication/sso_test",
		authenticationSSOTestHandler(),
	).
		Methods("GET").Name("authenticationBuildURLHandler")

	servicesRepo := repositories.NewServices(a.storage)
	servicesUseCase := usecases.NewServices(servicesRepo, serviceAccountsUseCase)
	authMiddle := authMiddleware(serviceAccountsUseCase)

	r.Handle(
		"/services",
		authMiddle(http.HandlerFunc(
			servicesCreateHandler(servicesUseCase),
		)),
	).
		Methods("POST").Name("servicesCreateHandler")

	r.Handle(
		"/service_accounts/{id}",
		authMiddle(http.HandlerFunc(
			serviceAccountsGetHandler(serviceAccountsUseCase),
		)),
	).
		Methods("GET").Name("serviceAccountsGetHandler")

	r.Handle(
		"/service_accounts",
		authMiddle(http.HandlerFunc(
			serviceAccountsCreateHandler(serviceAccountsUseCase),
		)),
	).
		Methods("POST").Name("serviceAccountsCreateHandler")

	rolesUseCase := usecases.NewRoles(rolesRepo, permissionsRepo)

	r.Handle(
		"/roles/{id}/permissions",
		authMiddle(http.HandlerFunc(
			rolesCreatePermissionHandler(serviceAccountsUseCase, rolesUseCase),
		)),
	).
		Methods("POST").Name("rolesCreatePermissionHandler")

	return r
}

//ListenAndServe requests
func (a *App) ListenAndServe() {
	listener, err := net.Listen("tcp", a.address)
	if err != nil {
		a.logger.WithError(err).Error("Failed to listen HTTP")
	}

	defer listener.Close()

	err = a.server.Serve(listener)
	if err != nil {
		a.logger.WithError(err).Error("Closed http listener")
	}
}
