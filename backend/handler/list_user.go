package handler

import (
	"context"
	"fmt"
	"restaurant-management-system/db"
	"restaurant-management-system/objects"
)

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
	UpdateAt  string `json:"updateAt"`
}

type ListUserRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ListUserResponse struct {
	objects.GenericResponse
	Users      []User `json:"users"`
	TotalCount int    `json:"totalCount"`
	Count      int    `json:"count"`
	Next       int    `json:"next"`
}

func ListUser(req objects.RequestObject) error {
	reqPay := req.Request.(*ListUserRequest)
	resPay := req.Response.(*ListUserResponse)

	if reqPay.Limit <= 0 {
		reqPay.Limit = 10
	}

	// Get total count
	var totalCount int
	err := db.Client.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&totalCount)
	if err != nil {
		return fmt.Errorf("failed to get total count: %w", err)
	}

	// Get users
	rows, err := db.Client.Query(context.Background(),
		"SELECT id, name, email, role, createdAt, updateAt FROM users LIMIT $1 OFFSET $2",
		reqPay.Limit, reqPay.Offset,
	)
	if err != nil {
		return fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		var createdAt, updateAt interface{}
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &createdAt, &updateAt)
		if err != nil {
			return fmt.Errorf("failed to scan user: %w", err)
		}
		if createdAt != nil {
			u.CreatedAt = fmt.Sprintf("%v", createdAt)
		}
		if updateAt != nil {
			u.UpdateAt = fmt.Sprintf("%v", updateAt)
		}
		users = append(users, u)
	}

	resPay.GenericResponse = objects.SuccessGeneric
	resPay.Users = users
	resPay.TotalCount = totalCount
	resPay.Count = len(users)

	if reqPay.Offset+len(users) < totalCount {
		resPay.Next = reqPay.Offset + reqPay.Limit
	} else {
		resPay.Next = 0
	}

	return nil
}
