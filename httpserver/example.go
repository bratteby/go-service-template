package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/bratteby/go-service-template/example"
)

type exampleHandler struct {
	exampleService exampleService
	encoder        encoder
}

func (h exampleHandler) GetRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", h.createExample)
		r.Get("/{id}", h.getExample)
	}
}

func (h *exampleHandler) createExample(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	var ex example.ExampleDTO
	if err := json.NewDecoder(r.Body).Decode(&ex); err != nil {
		h.encoder.respond(
			ctx, w,
			fmt.Errorf("could not decode request body %w", err),
			http.StatusBadRequest,
		)
		return
	}

	res, err := h.exampleService.CreateExample(ctx, ex)
	if err != nil {
		h.encoder.error(ctx, w, err)
		return
	}

	h.encoder.respond(ctx, w, res, http.StatusOK)
}

func (h *exampleHandler) getExample(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	exampleID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.encoder.respond(ctx, w, fmt.Errorf("could not parse ID from url"), http.StatusBadRequest)
		return
	}

	res, err := h.exampleService.GetExampleByID(ctx, exampleID)
	if err != nil {
		h.encoder.error(ctx, w, err)
		return
	}

	h.encoder.respond(ctx, w, res, http.StatusOK)
}
