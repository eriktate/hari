package http

import (
	"net/http"
	"strconv"

	"github.com/eriktate/hari"
	"github.com/go-chi/chi/v5"
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

func HandleGetWebhook(webhookService hari.WebhookService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := hari.ParseID(chi.URLParam(r, "id"))
		if err != nil {
			log.Error().Err(err).Msg("failed to parse Webhook ID")
			badRequest(w, hari.GetPublicError(hari.ErrCodeWebhookGetReq))
			return
		}

		webhook, err := webhookService.GetWebhook(ctx, id)
		if err != nil {
			log.Error().Err(err).Msg("failed to get Webhook")
			serverError(w, hari.GetPublicError(hari.ErrCodeWebhookGet))
			return
		}

		okJSON(w, webhook)
	}
}

func HandleQueryWebhook(webhookService hari.WebhookService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var query hari.WebhookQuery

		values := r.URL.Query()

		active := values.Get("active")
		name := values.Get("name")
		key := values.Get("key")

		if active != "" {
			activeBool, err := strconv.ParseBool(active)
			if err != nil {
				log.Error().Err(err).Str("active", active).Msg("invalid value for query param")
				badRequest(w, hari.GetPublicError(hari.ErrCodeWebhookQueryReq))
				return
			}

			query.Active = &activeBool
		}

		if name != "" {
			query.Name = &name
		}

		if key != "" {
			query.Key = &key
		}

		// TODO (etate): remove this hardcoded account ID once auth is working
		query.AccountID = hari.DefaultAccountID()
		webhooks, err := webhookService.QueryWebhooks(ctx, query)
		if err != nil {
			log.Error().Err(err).Msg("failed to query Webhook")
			serverError(w, hari.GetPublicError(hari.ErrCodeWebhookQuery))
			return
		}

		okJSON(w, webhooks)
	}
}

func HandleDeleteWebhook(webhookService hari.WebhookService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := hari.ParseID(chi.URLParam(r, "id"))
		if err != nil {
			log.Error().Err(err).Msg("failed to parse Webhook ID")
			badRequest(w, hari.GetPublicError(hari.ErrCodeWebhookDeleteReq))
			return
		}

		if err := webhookService.DeleteWebhook(ctx, id); err != nil {
			log.Error().Err(err).Msg("failed to get Webhook")
			serverError(w, hari.GetPublicError(hari.ErrCodeWebhookDelete))
			return
		}

		noContent(w)
	}
}
