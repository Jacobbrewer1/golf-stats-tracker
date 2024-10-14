package utils

import "context"

type contextKey string

const (
	userIdKey contextKey = "user_id"
)

// UserIdFromContext returns the user_id from the context.
func UserIdFromContext(ctx context.Context) int {
	userId, ok := ctx.Value(userIdKey).(int)
	if !ok {
		return -1
	}
	return userId
}

// UserIdToContext adds the user_id to the context.
func UserIdToContext(ctx context.Context, userId int) context.Context {
	return context.WithValue(ctx, userIdKey, userId)
}
