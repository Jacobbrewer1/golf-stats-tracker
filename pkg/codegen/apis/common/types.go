// Package common provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package common

// ErrorMessage defines the model for error_message.
type ErrorMessage struct {
	Error   *string `json:"error,omitempty"`
	Message *string `json:"message,omitempty"`
}

// Message defines the model for message.
type Message struct {
	Message *string `json:"message,omitempty"`
}

// LastId defines the model for last_id.
type LastId = string

// LastValue defines the model for last_value.
type LastValue = string

// LimitParam defines the model for limit_param.
type LimitParam = string

// SortBy defines the model for sort_by.
type SortBy = string

// SortDirection defines the model for sort_direction.
type SortDirection = string

// List of SortDirection
const (
	SortDirection_ASC  SortDirection = "ASC"
	SortDirection_DESC SortDirection = "DESC"
)
