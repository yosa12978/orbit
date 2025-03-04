package endpoints

import (
	"log/slog"
	"net/http"
	"orbit-app/internal/services"
	"orbit-app/pkg"
)

func GetSnippetByID(snippetService services.SnippetService, logger *slog.Logger) Endpoint {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := pkg.DecodeBase62(r.PathValue("id"))
		if err != nil {
			return EndpointError{
				StatusCode: http.StatusBadRequest,
				Err:        "invalid id",
			}
		}
		resp, err := snippetService.GetByID(r.Context(), id)
		if err != nil {
			logger.Error("error while fetching snippet by id", "error", err.Error())
			return EndpointError{
				StatusCode: http.StatusNotFound,
				Err:        "not found",
			}
		}
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte(resp.Content))
		return nil
	}
}
