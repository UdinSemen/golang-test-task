package numerator

import (
	"log"
	"net/http"

	httpAdapter "github.com/UdinSemen/golang-test-task/internal/adapters/http"

	"github.com/gin-gonic/gin"
)

type failResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message"`
}

type successResponse struct {
	Success bool    `json:"success" example:"true"`
	Message *string `json:"message,omitempty"`
	Body    any     `json:"body,omitempty"`
}

// AddNum godoc
//
//	@Summary	Добавление чисел
//
//	@Tags		API
//
//	@Accept		json
//
//	@Produce	json
//
//	@Param		request	body		http.AddNumRequest	true	"Структура добавления числа"
//
//	@Success	200		{object}	successResponse{body=http.NumResponse}
//
//	@Failure	422		{object}	failResponse
//
//	@Failure	500		{object}	failResponse
//
//	@Router		/num [post]
func (r *Router) AddNum(c *gin.Context) {
	var numRequest httpAdapter.AddNumRequest
	if err := c.ShouldBindJSON(&numRequest); err != nil {
		log.Println("failed to bind json", err)
		c.JSON(http.StatusUnprocessableEntity, failResponse{Success: false, Message: "invalid request"})
		return
	}

	nums, err := r.numeratorUsecase.AddNum(c.Request.Context(), numRequest.Num)
	if err != nil {
		log.Println("failed to add num", err)
		c.JSON(http.StatusInternalServerError, failResponse{Success: false, Message: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, successResponse{Success: true, Body: httpAdapter.NumResponse{Nums: nums}})
}
