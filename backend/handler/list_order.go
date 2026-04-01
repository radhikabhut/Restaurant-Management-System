package handler

import (
	"context"
	"fmt"
	"restaurant-management-system/db"
	"restaurant-management-system/objects"
	"time"

	"github.com/google/uuid"
)

type ListOrderRequest struct {
	UserId uuid.UUID `json:"userId"`
	Status string    `json:"status"`
}

type OrderItemResponse struct {
	Id       uuid.UUID `json:"id"`
	MenuId   uuid.UUID `json:"menuId"`
	MenuName string    `json:"menuName"`
	Quantity int       `json:"quantity"`
	Price    float64   `json:"price"`
}

type OrderResponse struct {
	Id          uuid.UUID           `json:"id"`
	UserId      uuid.UUID           `json:"userId"`
	TableNumber string              `json:"tableNumber"`
	TotalAmount float64             `json:"totalAmount"`
	Status      string              `json:"status"`
	IsDeleted   bool                `json:"isDeleted"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Items       []OrderItemResponse `json:"items"`
}

type ListOrderResponse struct {
	Orders []OrderResponse `json:"orders"`
	objects.GenericResponse
}

func ListOrders(req objects.RequestObject) error {
	reqPay := req.Request.(*ListOrderRequest)
	resPay := req.Response.(*ListOrderResponse)

	query := "SELECT id, userId, tableNumber, totalAmount, status, is_deleted, createdAt, updatedAt FROM orders WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if reqPay.UserId != uuid.Nil {
		query += fmt.Sprintf(" AND userId = $%d", argCount)
		args = append(args, reqPay.UserId)
		argCount++
	}
	if reqPay.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, reqPay.Status)
		argCount++
	}

	query += " ORDER BY createdAt DESC"

	rows, err := db.Client.Query(context.Background(), query, args...)
	if err != nil {
		return fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []OrderResponse
	for rows.Next() {
		var o OrderResponse
		var tableNum *string
		err := rows.Scan(&o.Id, &o.UserId, &tableNum, &o.TotalAmount, &o.Status, &o.IsDeleted, &o.CreatedAt, &o.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to scan order: %w", err)
		}
		if tableNum != nil {
			o.TableNumber = *tableNum
		}
		orders = append(orders, o)
	}

	// Fetch items for each order
	for i := range orders {
		itemRows, err := db.Client.Query(context.Background(),
			"SELECT oi.id, oi.menuId, m.name, oi.quantity, oi.price FROM order_items oi JOIN menus m ON oi.menuId = m.id WHERE oi.orderId = $1",
			orders[i].Id,
		)
		if err != nil {
			return fmt.Errorf("failed to query order items: %w", err)
		}
		defer itemRows.Close()

		var items []OrderItemResponse
		for itemRows.Next() {
			var it OrderItemResponse
			err := itemRows.Scan(&it.Id, &it.MenuId, &it.MenuName, &it.Quantity, &it.Price)
			if err != nil {
				return fmt.Errorf("failed to scan order item: %w", err)
			}
			items = append(items, it)
		}
		orders[i].Items = items
	}

	resPay.Orders = orders
	resPay.GenericResponse = objects.SuccessGeneric

	return nil
}
