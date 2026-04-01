package handler

import (
	"context"
	"fmt"
	"restaurant-management-system/db"
	"restaurant-management-system/objects"
	"restaurant-management-system/utils"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	Role   string `json:"role"`
	UserId string `json:"userId"`
	objects.GenericResponse
}

func Login(req objects.RequestObject) error {
	reqPay := req.Request.(*LoginRequest)
	resPay := req.Response.(*LoginResponse)

	var userID string
	var hashedPassword string
	var role string

	err := db.Client.QueryRow(context.Background(),
		"SELECT id, password, role FROM users WHERE email = $1",
		reqPay.Email,
	).Scan(&userID, &hashedPassword, &role)

	if err != nil {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "invalid email or password",
			StatusCode: 401,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(reqPay.Password))
	if err != nil {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "invalid email or password",
			StatusCode: 401,
		}
	}

	token, err := utils.GenerateToken(userID, role)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	resPay.Token = token
	resPay.Role = role
	resPay.UserId = userID
	resPay.GenericResponse = objects.SuccessGeneric

	return nil
}
