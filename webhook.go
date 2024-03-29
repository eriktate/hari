package hari

import (
	"context"
	"time"
)

// A Webhook holds all of the configuration required to run a webhook
type Webhook struct {
	ID             ID       `json:"id"`
	AccountID      ID       `json:"accountId"`
	Name           string   `json:"name"`
	Key            string   `json:"key"`
	Targets        []Target `json:"targets"`
	DefaultPayload *string  `json:"defaultPayload,omitEmpty"`
	Active         bool     `json:"active"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type NewWebhook struct {
	AccountID      ID      `json:"accountId"`
	Name           string  `json:"name"`
	Key            string  `json:"key"`
	DefaultPayload *string `json:"defaultPayload"`
}

type WebhookUpdates struct {
	Name           *string `json:"name"`
	DefaultPayload *string `json:"defaultPayload"`
	Active         *bool   `json:"active"`
}

type WebhookQuery struct {
	Active *bool   `json:"active"`
	Name   *string `json:"name"`
	Key    *string `json:"key"`
}

type WebhookService interface {
	CreateWebhook(ctx context.Context, webhook NewWebhook) (ID, error)
	UpdateWebhook(ctx context.Context, updates WebhookUpdates) error
	GetWebhook(ctx context.Context, id ID) (Webhook, error)
	QueryWebhooks(ctx context.Context, query WebhookQuery) ([]Webhook, error)
	DeleteWebhook(ctx context.Context, id ID) error
}
