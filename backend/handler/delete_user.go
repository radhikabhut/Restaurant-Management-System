package handler

import (
	"context"
	"fmt"
	"restaurant-management-system/db"
	"restaurant-management-system/objects"
)

type DeleteUserRequest struct {
	ID string `json:"id"`
}

type DeleteUserResponse struct {
	objects.GenericResponse
}

func DeleteUser(req objects.RequestObject) error {
	reqPay := req.Request.(*DeleteUserRequest)
	resPay := req.Response.(*DeleteUserResponse)

	// Validation
	if reqPay.ID == "" {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "id is required",
			StatusCode: 400,
		}
	}

	commandTag, err := db.Client.Exec(context.Background(), "DELETE FROM users WHERE id = $1", reqPay.ID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "user not found",
			StatusCode: 404,
		}
	}

	resPay.GenericResponse = objects.SuccessGeneric

	return nil
}
