package handler

import (
	"context"
	"fmt"
	"restaurant-management-system/db"
	"restaurant-management-system/objects"
)

type Menu struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Price       int64  `json:"price"`
	Category    string `json:"category"`
	IsAvailable bool   `json:"isAvailable"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type ListMenuRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ListMenuResponse struct {
	objects.GenericResponse
	Menus      []Menu `json:"menus"`
	TotalCount int    `json:"totalCount"`
	Count      int    `json:"count"`
	Next       int    `json:"next"`
}

func ListMenu(req objects.RequestObject) error {
	reqPay := req.Request.(*ListMenuRequest)
	resPay := req.Response.(*ListMenuResponse)

	if reqPay.Limit <= 0 {
		reqPay.Limit = 10
	}

	// Get total count
	var totalCount int
	err := db.Client.QueryRow(context.Background(), "SELECT COUNT(*) FROM menus").Scan(&totalCount)
	if err != nil {
		return fmt.Errorf("failed to get total count: %w", err)
	}

	// Get menus
	rows, err := db.Client.Query(context.Background(),
		"SELECT id, name, price, category, isAvailable, createdAt, updatedAt FROM menus LIMIT $1 OFFSET $2",
		reqPay.Limit, reqPay.Offset,
	)
	if err != nil {
		return fmt.Errorf("failed to query menus: %w", err)
	}
	defer rows.Close()

	var menus []Menu
	for rows.Next() {
		var m Menu
		var createdAt, updatedAt interface{}
		err := rows.Scan(&m.ID, &m.Name, &m.Price, &m.Category, &m.IsAvailable, &createdAt, &updatedAt)
		if err != nil {
			return fmt.Errorf("failed to scan menu: %w", err)
		}
		if createdAt != nil {
			m.CreatedAt = fmt.Sprintf("%v", createdAt)
		}
		if updatedAt != nil {
			m.UpdatedAt = fmt.Sprintf("%v", updatedAt)
		}
		menus = append(menus, m)
	}

	resPay.GenericResponse = objects.SuccessGeneric
	resPay.Menus = menus
	resPay.TotalCount = totalCount
	resPay.Count = len(menus)

	if reqPay.Offset+len(menus) < totalCount {
		resPay.Next = reqPay.Offset + reqPay.Limit
	} else {
		resPay.Next = 0
	}

	return nil
}
