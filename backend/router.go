package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"restaurant-management-system/middleware"
	"restaurant-management-system/view"
)

func defineRoutes(v1 *gin.RouterGroup) {
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Public/Login
	v1.POST("/login", view.Login)

	// Authenticated Group
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		// Menu Listing (All Roles)
		auth.POST("/menus/list", middleware.RBACMiddleware("menu::list"), view.ListMenu)

		// Menu Management (Admin Only)
		auth.POST("/menus/create", middleware.RBACMiddleware("menu::create"), view.CreateMenu)
		auth.POST("/menus/update", middleware.RBACMiddleware("menu::update"), view.UpdateMenu)
		auth.POST("/menus/delete", middleware.RBACMiddleware("menu::delete"), view.DeleteMenu)

		// User Management (Admin Only)
		auth.POST("/users/create", middleware.RBACMiddleware("user::create"), view.CreateUser)
		auth.POST("/users/list", middleware.RBACMiddleware("user::list"), view.ListUser)
		auth.POST("/users/update", middleware.RBACMiddleware("user::update"), view.UpdateUser)
		auth.POST("/users/delete", middleware.RBACMiddleware("user::delete"), view.DeleteUser)

		// Order Listing (All Roles)
		auth.POST("/orders/list", middleware.RBACMiddleware("order::list"), view.ListOrders)

		// Order Creation (Admin and Waiter)
		auth.POST("/orders/create", middleware.RBACMiddleware("order::create"), view.CreateOrder)

		// Order Updates (Admin and Kitchen)
		auth.POST("/orders/update", middleware.RBACMiddleware("order::update"), view.UpdateOrder)

		// Order Deletion (Admin Only)
		auth.POST("/orders/delete", middleware.RBACMiddleware("order::delete"), view.DeleteOrder)

		// Bill Generation (Admin and Waiter)
		auth.POST("/orders/bill", middleware.RBACMiddleware("order::bill"), view.GenerateBill)
	}
}
