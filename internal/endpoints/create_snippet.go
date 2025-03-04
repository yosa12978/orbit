package endpoints

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"orbit-app/internal/dto"
	"orbit-app/internal/services"
	"orbit-app/pkg"
)

const maxBodySize = 1 << 20

func CreateSnippet(
	snippetService services.SnippetService,
	logger *slog.Logger,
	hostname string,
) Endpoint {
	return func(w http.ResponseWriter, r *http.Request) error {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)
		if err != nil {
			if err == http.ErrBodyReadAfterClose {
				return EndpointError{
					http.StatusRequestEntityTooLarge,
					"Request body is too large",
				}
			}
			return EndpointError{
				StatusCode: http.StatusBadRequest,
				Err:        err.Error(),
			}
		}
		req := dto.SnippetCreateRequest{Content: string(data)}
		resp, err := snippetService.Create(r.Context(), req)
		if err != nil {
			if errors.As(err, &pkg.ValidationError{}) {
				return EndpointError{
					StatusCode: http.StatusBadRequest,
					Err:        err.Error(),
				}
			}
			logger.Error("error while creating snippet", "error", err.Error())
			return EndpointError{
				StatusCode: http.StatusInternalServerError,
				Err:        err.Error(),
			}
		}

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		fmt.Fprintf(w, "%s://%s/%v", scheme, hostname, pkg.EncodeBase62(resp.ID))
		return nil
	}
}
