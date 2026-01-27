package auth

import "context"

func Guard(ctx context.Context, roles ...string) error {
	return RequireRole(roles...)(ctx)
}
