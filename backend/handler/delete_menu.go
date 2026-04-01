package handler

import (
	"context"
	"fmt"
	"restaurant-management-system/db"
	"restaurant-management-system/objects"
)

type DeleteMenuRequest struct {
	ID string `json:"id"`
}

type DeleteMenuResponse struct {
	objects.GenericResponse
}

func DeleteMenu(req objects.RequestObject) error {
	reqPay := req.Request.(*DeleteMenuRequest)
	resPay := req.Response.(*DeleteMenuResponse)

	// Validation
	if reqPay.ID == "" {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "id is required",
			StatusCode: 400,
		}
	}

	commandTag, err := db.Client.Exec(context.Background(), "DELETE FROM menus WHERE id = $1", reqPay.ID)
	if err != nil {
		return fmt.Errorf("failed to delete menu: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "menu not found",
			StatusCode: 404,
		}
	}

	resPay.GenericResponse = objects.SuccessGeneric

	return nil
}
