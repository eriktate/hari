package hari

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type TargetStatus string

const (
	TargetStatusOk            = TargetStatus("ok")
	TargetStatusNotResponding = TargetStatus("not_responding")
	TargetStatusError         = TargetStatus("error")
	TargetStatusInvalid       = TargetStatus("invalid")
)

func ParseTargetStatus(status string) (TargetStatus, error) {
	possibleStatus := TargetStatus(strings.ToLower(status))
	switch possibleStatus {
	case TargetStatusOk, TargetStatusNotResponding, TargetStatusError:
		return possibleStatus, nil
	default:
		return TargetStatusInvalid, fmt.Errorf("invalid TargetStatus '%s'", status)
	}
}

type Target struct {
	ID        ID           `json:"id"`
	WebhookID ID           `json:"webhookId"`
	URL       url.URL      `json:"url"`
	Status    TargetStatus `json:"status"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

type HookStatus string

const (
	HookStatusSuccess   = HookStatus("success")   // successful hook, no issues
	HookStatusPending   = HookStatus("pending")   // currently hitting all registered targets
	HookStatusScheduled = HookStatus("scheduled") // hook is scheduled for the future
	HookStatusFailure   = HookStatus("failure")   // hook failed to reach one or more targets
	HookStatusInvalid   = HookStatus("invalid")   // parsing failure
)

func ParseHookStatus(status string) (HookStatus, error) {
	possibleStatus := HookStatus(strings.ToLower(status))
	switch possibleStatus {
	case HookStatusSuccess, HookStatusPending, HookStatusScheduled, HookStatusFailure:
		return possibleStatus, nil
	default:
		return HookStatusInvalid, fmt.Errorf("invalid HookStatus '%s'", status)
	}
}

// A Hook is an individual invocation of a Webhook. Hooks can be immediately invoked (most common) or even scheduled for the future
type Hook struct {
	ID        ID         `json:"id"`
	WebhookID ID         `json:"webhookId"`
	Status    HookStatus `json:"status"`
	Payload   *[]byte    `json:"payload"`
	RunAt     *time.Time `json:"runAt"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
