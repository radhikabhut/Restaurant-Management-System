package handler

import (
	"context"
	"fmt"
	"restaurant-management-system/db"
	"restaurant-management-system/objects"

	"github.com/google/uuid"
)

type DeleteOrderRequest struct {
	UserId uuid.UUID `json:"userId"`
	Id     uuid.UUID `json:"id"`
}

type DeleteOrderResponse struct {
	objects.GenericResponse
}

func DeleteOrder(req objects.RequestObject) error {
	reqPay := req.Request.(*DeleteOrderRequest)
	resPay := req.Response.(*DeleteOrderResponse)

	if reqPay.Id == uuid.Nil {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "id is required",
			StatusCode: 400,
		}
	}
	if reqPay.UserId == uuid.Nil {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "userId is required",
			StatusCode: 400,
		}
	}


	res, err := db.Client.Exec(context.Background(),
		"UPDATE orders SET is_deleted = true WHERE id = $1",
		reqPay.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if res.RowsAffected() == 0 {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "order not found",
			StatusCode: 404,
		}
	}

	resPay.GenericResponse = objects.SuccessGeneric

	return nil
}
