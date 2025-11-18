package numerator

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	AddNum(ctx context.Context, num int64) ([]int64, error)
}
type Router struct {
	numeratorUsecase Usecase
}

func NewRouter(numeratorUsecase Usecase) *Router {
	return &Router{
		numeratorUsecase: numeratorUsecase,
	}
}

func (r *Router) RegisterRoutes(ginGroup *gin.RouterGroup) {
	ginGroup.POST("/add", r.AddNum)
}
