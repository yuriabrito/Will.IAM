package testing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ghostec/Will.IAM/api"
	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/oauth2"
	"github.com/ghostec/Will.IAM/repositories"
	"github.com/ghostec/Will.IAM/usecases"
	"github.com/ghostec/Will.IAM/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GetConfig gets config for tests
func GetConfig(t *testing.T, path ...string) *viper.Viper {
	t.Helper()
	filePath := "./../testing/config.yaml"
	if len(path) > 0 {
		filePath = path[0]
	}
	config, err := utils.GetConfig(filePath)
	if err != nil {
		t.Fatal(err)
	}
	return config
}

// GetLogger gets config for tests
func GetLogger(t *testing.T) logrus.FieldLogger {
	t.Helper()
	return utils.GetLogger("0.0.0.0", 8080, 0, true)
}

// GetApp is a helper to create an *api.App
func GetApp(t *testing.T) *api.App {
	app, err := api.NewApp("0.0.0.0", 8080, GetConfig(t), GetLogger(t), nil)
	app.SetOAuth2Provider(oauth2.NewProviderBlankMock())
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return app
}

// DoRequest executes req over handler and returns a recorder
func DoRequest(
	t *testing.T, req *http.Request, handler http.Handler,
) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

// GetStorage returns a *repositories.Storage
func GetStorage(t *testing.T) *repositories.Storage {
	t.Helper()
	s := repositories.NewStorage()
	err := s.ConfigurePG(GetConfig(t))
	if err != nil {
		panic(err)
	}
	return s
}

// GetServiceAccountsUseCase returns a usecases.ServiceAccounts
func GetServiceAccountsUseCase(t *testing.T) usecases.ServiceAccounts {
	t.Helper()
	storage := GetStorage(t)
	saRepo := repositories.NewServiceAccounts(storage)
	rRepo := repositories.NewRoles(storage)
	pRepo := repositories.NewPermissions(storage)
	providerBlankMock := oauth2.NewProviderBlankMock()
	return usecases.NewServiceAccounts(saRepo, rRepo, pRepo, providerBlankMock)
}

// CreateRootServiceAccount creates a root service account with root access
func CreateRootServiceAccount(t *testing.T) *models.ServiceAccount {
	saUC := GetServiceAccountsUseCase(t)
	rootSA, err := saUC.CreateKeyPairType("root")
	if err != nil {
		panic(err)
	}
	return rootSA
}
