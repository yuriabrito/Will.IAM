package api

import (
	"encoding/json"
	"net/http"

	"github.com/ghostec/Will.IAM/usecases"
)

func amListHandler(
	amUC usecases.AM,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// l := middleware.GetLogger(r.Context())
		qs := r.URL.Query()
		prefixSl := qs["prefix"]
		prefix := ""
		if len(prefixSl) != 0 {
			prefix = prefixSl[0]
		}
		results, err := amUC.List(prefix)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		bts, err := json.Marshal(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		WriteBytes(w, http.StatusOK, bts)
	}
}
