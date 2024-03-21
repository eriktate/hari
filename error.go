package hari

type ErrCode int

const (
	ErrCodeWebhookCreate = iota
	ErrCodeWebhookCreateReq
	ErrCodeWebhookUpdate
	ErrCodeWebhookUpdateReq
)

// HariErrors are meant to be safe for public exposure. In most cases these are pre-canned messages and don't expose actual error details
type HariError struct {
	Code ErrCode `json:"code"`
	Msg  string  `json:"msg"`
}

// NOTE: the order that these errors show up _must_ match the order of the ErrCode "enum" above, otherwise lookups will fail or become incorrect
var publicErrors = [...]HariError{
	{
		Code: ErrCodeWebhookCreate,
		Msg:  "Looks like your request was correct but we weren't able to create your webhook. Please wait a few seconds and try again",
	},
	{
		Code: ErrCodeWebhookCreateReq,
		Msg:  "Looks like your request might have been incorrect. Make sure you provide all required fields and try again",
	},
	{
		Code: ErrCodeWebhookUpdate,
		Msg:  "Looks like we weren't able to update your webhook. Please wait a few seconds and try again",
	},
	{
		Code: ErrCodeWebhookUpdateReq,
		Msg:  "Looks like your request might have been incorrect. Make sure you provide all required fields and try again",
	},
}

func GetPublicError(code ErrCode) HariError {
	return publicErrors[code]
}
