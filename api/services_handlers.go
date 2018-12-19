package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/usecases"
)

func servicesCreateHandler(
	servicesUseCase usecases.Services,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			body, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				return err
			}
			service := &models.Service{}
			err = json.Unmarshal(body, service)
			if err != nil {
				return err
			}
			// TODO: check if user has William::RL::CreateService::*
			saID, ok := getServiceAccountID(r.Context())
			if !ok {
				return fmt.Errorf("service_account_id not set in ctx")
			}
			if err := servicesUseCase.Create(service, saID); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
