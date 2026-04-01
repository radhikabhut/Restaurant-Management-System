package handler

import (
	"context"
	"fmt"
	"restaurant-management-system/db"
	"restaurant-management-system/objects"
	"time"

	"github.com/google/uuid"
)

type OrderItemRequest struct {
	MenuId   uuid.UUID `json:"menuId"`
	Quantity int       `json:"quantity"`
}

type CreateOrderRequest struct {
	UserId      uuid.UUID          `json:"userId"`
	TableNumber string             `json:"tableNumber"`
	OrderItems  []OrderItemRequest `json:"orderItems"`
}

type CreateOrderResponse struct {
	Id         uuid.UUID `json:"id"`
	objects.GenericResponse
}

func CreateOrder(req objects.RequestObject) error {
	reqPay := req.Request.(*CreateOrderRequest)
	resPay := req.Response.(*CreateOrderResponse)

	// Validation
	if reqPay.UserId == uuid.Nil {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "userId is required",
			StatusCode: 400,
		}
	}
	if len(reqPay.OrderItems) == 0 {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "orderItems are required",
			StatusCode: 400,
		}
	}

	// Start a transaction
	tx, err := db.Client.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	// Generate Order ID (UUID v7)
	orderID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate order uuid: %w", err)
	}

	var totalAmount float64
	now := time.Now()

	type itemDetail struct {
		id       uuid.UUID
		menuId   uuid.UUID
		quantity int
		price    float64
	}
	var itemsToInsert []itemDetail

	// 1. Calculate prices and total amount first
	for _, item := range reqPay.OrderItems {
		if item.Quantity <= 0 {
			return objects.GenericError{
				Success:    false,
				ErrMsg:     "quantity must be greater than 0",
				StatusCode: 400,
			}
		}

		var price float64
		err = tx.QueryRow(context.Background(),
			"SELECT price FROM menus WHERE id = $1",
			item.MenuId,
		).Scan(&price)
		if err != nil {
			return objects.GenericError{
				Success:    false,
				ErrMsg:     fmt.Sprintf("menu item %s not found", item.MenuId),
				StatusCode: 404,
			}
		}

		itemTotal := price * float64(item.Quantity)
		totalAmount += itemTotal

		orderItemID, _ := uuid.NewV7()
		itemsToInsert = append(itemsToInsert, itemDetail{
			id:       orderItemID,
			menuId:   item.MenuId,
			quantity: item.Quantity,
			price:    price,
		})
	}

	// 2. Insert into orders table (Parent first)
	_, err = tx.Exec(context.Background(),
		"INSERT INTO orders (id, userId, totalAmount, status, tableNumber, createdAt, updatedAt) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		orderID, reqPay.UserId, totalAmount, "pending", reqPay.TableNumber, now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	// 3. Insert into order_items table (Children)
	for _, item := range itemsToInsert {
		_, err = tx.Exec(context.Background(),
			"INSERT INTO order_items (id, orderId, menuId, quantity, price, createdAt, updatedAt) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			item.id, orderID, item.menuId, item.quantity, item.price, now, now,
		)
		if err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	resPay.Id = orderID
	resPay.GenericResponse = objects.SuccessGeneric

	return nil
}
