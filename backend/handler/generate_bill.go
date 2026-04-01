package handler

import (
	"context"
	"fmt"
	"restaurant-management-system/db"
	"restaurant-management-system/objects"
	"time"

	"github.com/google/uuid"
)

type GenerateBillRequest struct {
	OrderId uuid.UUID `json:"orderId"`
}

type BillItem struct {
	MenuName string  `json:"menuName"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Subtotal float64 `json:"subtotal"`
}

type GenerateBillResponse struct {
	objects.GenericResponse
	OrderId       uuid.UUID  `json:"orderId"`
	TableNumber   string     `json:"tableNumber"`
	Items         []BillItem `json:"items"`
	TotalAmount   float64    `json:"totalAmount"`    // Sum of items
	Tax           float64    `json:"tax"`            // 5%
	ServiceCharge float64    `json:"serviceCharge"`  // 10%
	GrandTotal    float64    `json:"grandTotal"`
	BillDate      time.Time  `json:"billDate"`
}

func GenerateBill(req objects.RequestObject) error {
	reqPay := req.Request.(*GenerateBillRequest)
	resPay := req.Response.(*GenerateBillResponse)

	if reqPay.OrderId == uuid.Nil {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "orderId is required",
			StatusCode: 400,
		}
	}

	// Fetch Order and Table Number
	var tableNumberStr string
	var tableNumber *string
	var totalAmount float64
	err := db.Client.QueryRow(context.Background(),
		"SELECT tableNumber, totalAmount FROM orders WHERE id = $1",
		reqPay.OrderId,
	).Scan(&tableNumber, &totalAmount)
	if err != nil {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "order not found",
			StatusCode: 404,
		}
	}
	if tableNumber != nil {
		tableNumberStr = *tableNumber
	}

	// Fetch Items
	rows, err := db.Client.Query(context.Background(),
		"SELECT m.name, oi.quantity, oi.price FROM order_items oi JOIN menus m ON oi.menuId = m.id WHERE oi.orderId = $1",
		reqPay.OrderId,
	)
	if err != nil {
		return fmt.Errorf("failed to query order items: %w", err)
	}
	defer rows.Close()

	var items []BillItem
	for rows.Next() {
		var bi BillItem
		err := rows.Scan(&bi.MenuName, &bi.Quantity, &bi.Price)
		if err != nil {
			return fmt.Errorf("failed to scan bill item: %w", err)
		}
		bi.Subtotal = bi.Price * float64(bi.Quantity)
		items = append(items, bi)
	}

	// Calculations
	tax := totalAmount * 0.05
	serviceCharge := totalAmount * 0.10
	grandTotal := totalAmount + tax + serviceCharge

	// Response
	resPay.GenericResponse = objects.SuccessGeneric
	resPay.OrderId = reqPay.OrderId
	resPay.TableNumber = tableNumberStr
	resPay.Items = items
	resPay.TotalAmount = totalAmount
	resPay.Tax = tax
	resPay.ServiceCharge = serviceCharge
	resPay.GrandTotal = grandTotal
	resPay.BillDate = time.Now()

	return nil
}
