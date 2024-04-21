package http

import (
	"encoding/json"
	"net/http"

	"github.com/eriktate/hari"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Addr           string
	WebhookService hari.WebhookService
}

func getRoutes(cfg Config) http.Handler {
	r := chi.NewRouter()

	// webhook routes
	r.Post("/webhook", HandleCreateWebhook(cfg.WebhookService))
	r.Put("/webhook", HandleUpdateWebhook(cfg.WebhookService))
	r.Get("/webhook/{id}", HandleGetWebhook(cfg.WebhookService))
	r.Get("/webhook", HandleQueryWebhook(cfg.WebhookService))
	r.Delete("/webhook/{id}", HandleDeleteWebhook(cfg.WebhookService))

	return r
}

func Serve(cfg Config) error {
	return http.ListenAndServe(cfg.Addr, getRoutes(cfg))
}

// response writers
func respond(w http.ResponseWriter, payload []byte) {
	if _, err := w.Write(payload); err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}

func respondWithJSON(w http.ResponseWriter, payload any, status int) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal payload into response")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	respond(w, data)
}

func okJSON(w http.ResponseWriter, payload any) {
	respondWithJSON(w, payload, http.StatusOK)
}

func createdJSON(w http.ResponseWriter, payload any) {
	respondWithJSON(w, payload, http.StatusCreated)
}

func noContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	respond(w, nil)
}

func badRequest(w http.ResponseWriter, err hari.HariError) {
	respondWithJSON(w, err, http.StatusBadRequest)
}

func serverError(w http.ResponseWriter, err hari.HariError) {
	respondWithJSON(w, err, http.StatusInternalServerError)
}
