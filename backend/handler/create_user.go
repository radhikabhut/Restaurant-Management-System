package handler

import (
	"context"
	"fmt"
	"restaurant-management-system/db"
	"restaurant-management-system/objects"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CreateUserResponse struct {
	objects.GenericResponse
}

func CreateUser(req objects.RequestObject) error {
	reqPay := req.Request.(*CreateUserRequest)
	resPay := req.Response.(*CreateUserResponse)

	// Validation
	if reqPay.Name == "" {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "name is required",
			StatusCode: 400,
		}
	}
	if reqPay.Email == "" {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "email is required",
			StatusCode: 400,
		}
	}
	if reqPay.Password == "" {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "password is required",
			StatusCode: 400,
		}
	}
	if reqPay.Role == "" {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "role is required",
			StatusCode: 400,
		}
	}

	// Password Hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqPay.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate UUID v7
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate uuid: %w", err)
	}

	now := time.Now()
	_, err = db.Client.Exec(context.Background(),
		"INSERT INTO users (id, name, email, password, role, createdAt, updateAt) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		id, reqPay.Name, reqPay.Email, string(hashedPassword), reqPay.Role, now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	resPay.GenericResponse = objects.SuccessGeneric

	return nil
}
