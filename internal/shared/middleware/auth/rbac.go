package auth

import (
	"context"

	"github.com/nassabiq/golang-template/internal/modules/auth/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RequireRole(allowedRoles ...string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		_, roleID, ok := FromContext(ctx)
		if !ok {
			return status.Error(codes.Unauthenticated, "unauthorized")
		}

		// Convert Role ID to Role Name
		roleName, exists := domain.RoleIDToName[domain.RoleID(roleID)]
		if !exists {
			// Jika ID tidak dikenal, return error
			return status.Error(codes.PermissionDenied, "invalid role id")
		}

		for _, r := range allowedRoles {
			if string(roleName) == r {
				return nil
			}
		}

		return status.Error(codes.PermissionDenied, "forbidden")
	}
}
