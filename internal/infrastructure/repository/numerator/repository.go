package numerator

import (
	"context"
	"fmt"

	"github.com/UdinSemen/golang-test-task/internal/domain/usecase/numerator"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ numerator.Repository = new(Repository)

type Repository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) *Repository {
	return &Repository{conn: conn}
}

func (r *Repository) AddNum(ctx context.Context, num int64) error {
	const op = "numerator.Repository.AddNum"

	_, err := r.conn.Exec(ctx, "INSERT INTO public.nums(num) VALUES($1)", num)
	if err != nil {
		return fmt.Errorf("%s.%s: %w", op, "conn.Exec", err)
	}

	return nil
}

func (r *Repository) GetSortedNums(ctx context.Context) ([]int64, error) {
	const op = "numerator.Repository.GetSortedNums"

	rows, err := r.conn.Query(ctx, `SELECT num FROM public.nums ORDER BY num`)
	if err != nil {
		return nil, fmt.Errorf("%s.%s: %w", op, "conn.Query", err)
	}
	defer rows.Close()

	nums := make([]int64, 0)
	for rows.Next() {
		var num int64
		if err = rows.Scan(&num); err != nil {
			return nil, fmt.Errorf("%s.%s: %w", op, "rows.Scan", err)
		}
		nums = append(nums, num)
	}

	return nums, nil
}
