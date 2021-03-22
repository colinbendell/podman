package libpod

import (
	"net/http"

	"github.com/containers/podman/v3/libpod"
	"github.com/containers/podman/v3/pkg/api/handlers/utils"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman/v3/pkg/domain/infra/abi"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

func CreateSecret(w http.ResponseWriter, r *http.Request) {
	var (
		runtime = r.Context().Value("runtime").(*libpod.Runtime)
		decoder = r.Context().Value("decoder").(*schema.Decoder)
	)
	query := struct {
		Name   string `schema:"name"`
		Driver string `schema:"driver"`
	}{
		// override any golang type defaults
	}
	opts := entities.SecretCreateOptions{}
	if err := decoder.Decode(&query, r.URL.Query()); err != nil {
		utils.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest,
			errors.Wrapf(err, "failed to parse parameters for %s", r.URL.String()))
		return
	}
	opts.Driver = query.Driver

	ic := abi.ContainerEngine{Libpod: runtime}
	report, err := ic.SecretCreate(r.Context(), query.Name, r.Body, opts)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
	utils.WriteResponse(w, http.StatusOK, report)
}