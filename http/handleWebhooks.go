package http

import (
	"net/http"

	"github.com/eriktate/hari"
	"github.com/rs/zerolog/log"
)

func HandleCreateWebhook(webhookService hari.WebhookService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		webhook, err := parseJSON[hari.NewWebhook](r)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse Webhook")
			badRequest(w, hari.GetPublicError(hari.ErrCodeWebhookCreateReq))
			return
		}

		id, err := webhookService.CreateWebhook(ctx, webhook)
		if err != nil {
			log.Error().Err(err).Msg("failed to create Webhook")
			serverError(w, hari.GetPublicError(hari.ErrCodeWebhookCreate))
			return
		}

		createdJSON(w, id)
	}
}

func HandleUpdateWebhook(webhookService hari.WebhookService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		updates, err := parseJSON[hari.WebhookUpdates](r)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse Webhook")
			badRequest(w, hari.GetPublicError(hari.ErrCodeWebhookUpdateReq))
			return
		}

		if err := webhookService.UpdateWebhook(ctx, updates); err != nil {
			log.Error().Err(err).Msg("failed to create Webhook")
			serverError(w, hari.GetPublicError(hari.ErrCodeWebhookUpdate))
			return
		}

		noContent(w)
	}
}
