package auth

import "context"

type (
	userIDKey struct{}
	roleKey   struct{}
)

func WithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}
func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey{}, role)
}

func UserIDFromCtx(ctx context.Context) int64 {
	if v := ctx.Value(userIDKey{}); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

func RoleFromCtx(ctx context.Context) string {
	if v := ctx.Value(roleKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
