package handler

import (
	"context"
	"fmt"
	"restaurant-management-system/db"
)

// CheckPermission checks if a role has a specific permission in the database.
func CheckPermission(ctx context.Context, role string, permission string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM role_permission WHERE role = $1 AND permission = $2)`
	err := db.Client.QueryRow(ctx, query, role, permission).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}
	return exists, nil
}
