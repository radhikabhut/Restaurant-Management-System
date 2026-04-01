package handler

import (
	"fmt"
	"restaurant-management-system/db"
	"restaurant-management-system/objects"
	"time"

	"github.com/google/uuid"
)

type UpdateOrderRequest struct {
	Id             uuid.UUID           `json:"id"`
	Status         string              `json:"status"`
	OrderItems     []OrderItemRequest `json:"orderItems"`
	DeletedItemIds []uuid.UUID         `json:"deletedItemIds"`
}

type UpdateOrderResponse struct {
	objects.GenericResponse
}

func UpdateOrder(req objects.RequestObject) error {
	reqPay := req.Request.(*UpdateOrderRequest)
	resPay := req.Response.(*UpdateOrderResponse)

	if reqPay.Id == uuid.Nil {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "id is required",
			StatusCode: 400,
		}
	}

	userRole, ok := req.Ctx.Value("userRole").(string)
	if !ok {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "user role missing in context",
			StatusCode: 500,
		}
	}

	// Permission logic:
	// 1. Status change: Admin or Kitchen
	// 2. Add items: Admin or Waiter

	if reqPay.Status != "" {
		if userRole != "kitchen" && userRole != "admin" {
			return objects.GenericError{
				Success:    false,
				ErrMsg:     "only kitchen staff or admin can change order status",
				StatusCode: 403,
			}
		}
	}

	if len(reqPay.OrderItems) > 0 || len(reqPay.DeletedItemIds) > 0 {
		if userRole != "waiter" && userRole != "admin" {
			return objects.GenericError{
				Success:    false,
				ErrMsg:     "only waiter or admin can add or delete items in order",
				StatusCode: 403,
			}
		}
	}

	// Check if order is deleted
	var isDeleted bool
	err := db.Client.QueryRow(req.Ctx, "SELECT is_deleted FROM orders WHERE id = $1", reqPay.Id).Scan(&isDeleted)
	if err != nil {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "order not found",
			StatusCode: 404,
		}
	}
	if isDeleted {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "cannot update a deleted order",
			StatusCode: 400,
		}
	}

	// Start transaction
	tx, err := db.Client.Begin(req.Ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(req.Ctx)

	now := time.Now()

	// Update Status if provided
	if reqPay.Status != "" {
		_, err = tx.Exec(req.Ctx,
			"UPDATE orders SET status = $1, updatedAt = $2 WHERE id = $3",
			reqPay.Status, now, reqPay.Id,
		)
		if err != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}
	}

	// Add Order Items if provided
	if len(reqPay.OrderItems) > 0 {
		var additionalAmount float64
		for _, item := range reqPay.OrderItems {
			if item.Quantity <= 0 {
				return objects.GenericError{
					Success:    false,
					ErrMsg:     "quantity must be greater than 0",
					StatusCode: 400,
				}
			}

			var price float64
			err = tx.QueryRow(req.Ctx,
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
			additionalAmount += itemTotal

			orderItemID, _ := uuid.NewV7()
			_, err = tx.Exec(req.Ctx,
				"INSERT INTO order_items (id, orderId, menuId, quantity, price, createdAt, updatedAt) VALUES ($1, $2, $3, $4, $5, $6, $7)",
				orderItemID, reqPay.Id, item.MenuId, item.Quantity, price, now, now,
			)
			if err != nil {
				return fmt.Errorf("failed to insert order item: %w", err)
			}
		}

		// Update totalAmount in orders
		_, err = tx.Exec(req.Ctx,
			"UPDATE orders SET totalAmount = totalAmount + $1, updatedAt = $2 WHERE id = $3",
			additionalAmount, now, reqPay.Id,
		)
		if err != nil {
			return fmt.Errorf("failed to update order amount: %w", err)
		}
	}

	// Delete Order Items if provided
	if len(reqPay.DeletedItemIds) > 0 {
		var totalToSubtract float64
		for _, itemID := range reqPay.DeletedItemIds {
			var price float64
			var quantity int
			err = tx.QueryRow(req.Ctx,
				"SELECT price, quantity FROM order_items WHERE id = $1 AND orderId = $2",
				itemID, reqPay.Id,
			).Scan(&price, &quantity)
			if err != nil {
				// Item might have been deleted already or doesn't belong to this order
				continue
			}

			totalToSubtract += price * float64(quantity)

			_, err = tx.Exec(req.Ctx, "DELETE FROM order_items WHERE id = $1", itemID)
			if err != nil {
				return fmt.Errorf("failed to delete order item: %w", err)
			}
		}

		// Update totalAmount in orders
		_, err = tx.Exec(req.Ctx,
			"UPDATE orders SET totalAmount = totalAmount - $1, updatedAt = $2 WHERE id = $3",
			totalToSubtract, now, reqPay.Id,
		)
		if err != nil {
			return fmt.Errorf("failed to update order amount after deletion: %w", err)
		}
	}

	if err := tx.Commit(req.Ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	resPay.GenericResponse = objects.SuccessGeneric

	return nil
}
