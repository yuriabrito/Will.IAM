package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/usecases"
	"github.com/topfreegames/extensions/middleware"
)

func servicesListHandler(
	ssUC usecases.Services,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
		ssSl, err := ssUC.WithCtx(r.Context()).List()
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		bts, err := keepJSONFieldsBytes(ssSl, "id", "name", "created_at", "updated_at")
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		WriteBytes(w, 200, bts)
	}
}

func servicesCreateHandler(
	ssUC usecases.Services,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
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
			if err := ssUC.WithCtx(r.Context()).Create(service, saID); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
