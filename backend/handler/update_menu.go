package handler

import (
	"context"
	"fmt"
	"restaurant-management-system/db"
	"restaurant-management-system/objects"
	"time"
)

type UpdateMenuRequest struct {
	ID          string   `json:"id"`
	Name        *string  `json:"name"`
	Price       *float64 `json:"price"`
	Category    *string  `json:"category"`
	IsAvailable *bool    `json:"isAvailable"`
}

type UpdateMenuResponse struct {
	objects.GenericResponse
}

func UpdateMenu(req objects.RequestObject) error {
	reqPay := req.Request.(*UpdateMenuRequest)
	resPay := req.Response.(*UpdateMenuResponse)

	// Validation
	if reqPay.ID == "" {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "id is required",
			StatusCode: 400,
		}
	}

	if reqPay.Price != nil && *reqPay.Price <= 0 {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "price must be greater than 0",
			StatusCode: 400,
		}
	}

	now := time.Now()
	commandTag, err := db.Client.Exec(context.Background(),
		`UPDATE menus SET 
			name = COALESCE($1, name), 
			price = COALESCE($2, price), 
			category = COALESCE($3, category), 
			isAvailable = COALESCE($4, isAvailable), 
			updatedAt = $5 
		WHERE id = $6`,
		reqPay.Name, reqPay.Price, reqPay.Category, reqPay.IsAvailable, now, reqPay.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update menu: %w", err)
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
