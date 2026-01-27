package auth

import "context"

type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

func WithUser(ctx context.Context, userID, role string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, roleKey, role)
	return ctx
}

func FromContext(ctx context.Context) (userID string, role string, ok bool) {
	userID, ok1 := ctx.Value(userIDKey).(string)
	role, ok2 := ctx.Value(roleKey).(string)

	if !ok1 || !ok2 {
		return "", "", false
	}

	return userID, role, true
}
