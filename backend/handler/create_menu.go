package handler

import (
	"context"
	"fmt"
	"restaurant-management-system/db"
	"restaurant-management-system/objects"
	"time"

	"github.com/google/uuid"
)

type CreateMenuRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	IsAvailable bool    `json:"isAvailable"`
}

type CreateMenuResponse struct {
	objects.GenericResponse
}

func CreateMenu(req objects.RequestObject) error {
	reqPay := req.Request.(*CreateMenuRequest)
	resPay := req.Response.(*CreateMenuResponse)

	// Validation
	if reqPay.Name == "" {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "name is required",
			StatusCode: 400,
		}
	}
	if reqPay.Price <= 0 {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "price must be greater than 0",
			StatusCode: 400,
		}
	}
	if reqPay.Category == "" {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "category is required",
			StatusCode: 400,
		}
	}

	// Generate UUID v7
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate uuid: %w", err)
	}

	now := time.Now()
	_, err = db.Client.Exec(context.Background(),
		"INSERT INTO menus (id, name, price, category, isAvailable, createdAt, updatedAt) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		id, reqPay.Name, reqPay.Price, reqPay.Category, reqPay.IsAvailable, now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to insert menu: %w", err)
	}

	resPay.GenericResponse = objects.SuccessGeneric

	return nil
}
