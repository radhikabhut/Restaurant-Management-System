package handler

import (
	"context"
	"fmt"
	"restaurant-management-system/db"
	"restaurant-management-system/objects"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UpdateUserRequest struct {
	ID       string  `json:"id"`
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Role     *string `json:"role"`
}

type UpdateUserResponse struct {
	objects.GenericResponse
}

func UpdateUser(req objects.RequestObject) error {
	reqPay := req.Request.(*UpdateUserRequest)
	resPay := req.Response.(*UpdateUserResponse)

	// Validation
	if reqPay.ID == "" {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "id is required",
			StatusCode: 400,
		}
	}

	var hashedPassword *string
	if reqPay.Password != nil && *reqPay.Password != "" {
		hp, err := bcrypt.GenerateFromPassword([]byte(*reqPay.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		s := string(hp)
		hashedPassword = &s
	}

	now := time.Now()
	commandTag, err := db.Client.Exec(context.Background(),
		`UPDATE users SET 
			name = COALESCE($1, name), 
			email = COALESCE($2, email), 
			password = COALESCE($3, password), 
			role = COALESCE($4, role), 
			updateAt = $5 
		WHERE id = $6`,
		reqPay.Name, reqPay.Email, hashedPassword, reqPay.Role, now, reqPay.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "user not found",
			StatusCode: 404,
		}
	}

	resPay.GenericResponse = objects.SuccessGeneric

	return nil
}
