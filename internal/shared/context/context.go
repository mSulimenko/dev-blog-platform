package context

import (
	"context"
	"fmt"
)

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

func GetUserRoleFromContext(ctx context.Context) (string, error) {
	role, ok := ctx.Value("user_role").(string)
	if !ok || role == "" {
		return "", fmt.Errorf("user role not found in context")
	}
	return role, nil
}
