package numerator

import (
	"context"
	"fmt"

	"github.com/UdinSemen/golang-test-task/internal/infrastructure/api/http/numerator"
)

type Repository interface {
	AddNum(ctx context.Context, num int64) error
	GetSortedNums(ctx context.Context) ([]int64, error)
}

var _ numerator.Usecase = new(Usecase)

type Usecase struct {
	numeratorRepository Repository
}

func NewUsecase(numeratorRepository Repository) *Usecase {
	return &Usecase{numeratorRepository: numeratorRepository}
}

func (u *Usecase) AddNum(
	ctx context.Context,
	num int64,
) ([]int64, error) {
	const op = "usecase.numerator.AddNum"

	if err := u.numeratorRepository.AddNum(ctx, num); err != nil {
		return nil, fmt.Errorf("%s.%s:%w", op, "AddNum", err)
	}

	nums, err := u.numeratorRepository.GetSortedNums(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s.%s:%w", op, "GetSortedNums", err)
	}

	return nums, nil
}
