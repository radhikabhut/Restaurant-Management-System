package view

import (
	"log/slog"
	"net/http"
	"restaurant-management-system/handler"
	"restaurant-management-system/objects"

	"github.com/gin-gonic/gin"
)

func ListOrders(ctx *gin.Context) {
	var req objects.RequestObject

	var reqPay handler.ListOrderRequest
	var resPay handler.ListOrderResponse

	if err := ctx.ShouldBindJSON(&reqPay); err != nil {
		slog.Error("error while binding json", "error", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, objects.MalformedError)
		return
	}

	req.Request = &reqPay
	req.Response = &resPay

	if err := handler.ListOrders(req); err != nil {
		slog.Error("error while listing orders", "error", err.Error())
		genericErr, ok := err.(objects.GenericError)
		if !ok {
			slog.Error("error while casting error to generic error: ", "error", err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, objects.InternalError)
			return
		}
		ctx.AbortWithStatusJSON(genericErr.StatusCode, genericErr)
		return
	}

	ctx.JSON(http.StatusOK, resPay)
}
