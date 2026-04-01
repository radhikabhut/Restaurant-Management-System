package view

import (
	"log/slog"
	"net/http"
	"restaurant-management-system/handler"
	"restaurant-management-system/objects"

	"github.com/gin-gonic/gin"
)

func DeleteMenu(ctx *gin.Context) {
	var req objects.RequestObject

	var reqPay handler.DeleteMenuRequest
	var resPay handler.DeleteMenuResponse

	if err := ctx.ShouldBindJSON(&reqPay); err != nil {
		slog.Error("error while binding json", "error", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, objects.MalformedError)
		return
	}

	req.Request = &reqPay
	req.Response = &resPay

	if err := handler.DeleteMenu(req); err != nil {
		slog.Error("error while deleting menu", "error", err.Error())
		err, ok := err.(objects.GenericError)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, objects.InternalError)
			return
		}
		ctx.AbortWithStatusJSON(err.StatusCode, err)
		return
	}

	ctx.JSON(http.StatusOK, resPay)
}
