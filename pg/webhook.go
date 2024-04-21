package pg

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/eriktate/hari"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/rs/zerolog/log"
)

func (pg PG) CreateWebhook(ctx context.Context, wh hari.NewWebhook) (hari.ID, error) {
	id := hari.NewID()

	query := `
	insert into webhooks (id, account_id, name, key)
	values ($1, $2, $3, $4)
	`

	if _, err := pg.pool.Exec(ctx, query, id, wh.AccountID, wh.Name, wh.Key); err != nil {
		return hari.NullID(), fmt.Errorf("failed to CreateWebhook: %w", err)
	}

	return id, nil
}

func (pg PG) UpdateWebhook(ctx context.Context, updates hari.WebhookUpdates) error {
	q := sq.Update("webhooks").Where("id = ?", updates.ID)
	if updates.Name != nil {
		q.Set("name", updates.Name)
	}

	if updates.DefaultPayload != nil {
		q.Set("default_payload", updates.DefaultPayload)
	}

	if updates.Active != nil {
		q.Set("active", updates.Active)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("failed to generate query for UpdateWebhook: %w", err)
	}

	if _, err := pg.pool.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to UpdateWebhook: %w", err)
	}

	return nil
}

func (pg PG) GetWebhook(ctx context.Context, id hari.ID) (hari.Webhook, error) {
	var webhook hari.Webhook
	query := "select * from webhooks where id = $1 and deleted_at is null"

	if err := pgxscan.Get(ctx, pg.pool, &webhook, query, id); err != nil {
		return webhook, fmt.Errorf("failed to GetWebhook: %w", err)
	}

	return webhook, nil
}

func (pg PG) QueryWebhooks(ctx context.Context, query hari.WebhookQuery) ([]hari.Webhook, error) {
	var webhooks []hari.Webhook

	q := sq.Select("*").From("webhooks").Where("account_id = $1 and deleted_at is null", query.AccountID)

	if query.Active != nil {
		q.Where("active = ?", query.Active)
	}

	if query.Name != nil {
		q.Where("name = ?", query.Name)
	}

	if query.Key != nil {
		q.Where("key = ?", query.Key)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate query QueryWebhooks: %w", err)
	}

	log.Info().Any("args", args).Msg(sql)
	if err := pgxscan.Select(ctx, pg.pool, &webhooks, sql, args...); err != nil {
		return nil, fmt.Errorf("failed to QueryWebhooks: %w", err)
	}

	return webhooks, nil
}

func (pg PG) DeleteWebhook(ctx context.Context, id hari.ID) error {
	query := "update webhooks set deleted_at = now() where id = $1"

	if _, err := pg.pool.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to DeleteWebhooks: %w", err)
	}

	return nil
}
