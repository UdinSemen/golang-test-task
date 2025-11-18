//go:build integration

package numerator_test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/UdinSemen/golang-test-task/internal/infrastructure/repository/numerator"
	testconnectors "github.com/UdinSemen/golang-test-task/pkg/testconnectors/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	postgresCont "github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	sqlDB *pgxpool.Pool
	repo  *numerator.Repository
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgresCont.Run(ctx,
		"postgres:16-alpine",
		postgresCont.WithDatabase(dbName),
		postgresCont.WithUsername(dbUser),
		postgresCont.WithPassword(dbPassword),
		postgresCont.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}
	defer func() {
		if errCont := testcontainers.TerminateContainer(postgresContainer); errCont != nil {
			log.Printf("failed to terminate container: %v", errCont)
		}
	}()

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get container host: %v", err)
	}
	portRaw, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("failed to get mapped port: %v", err)
	}
	port, _ := strconv.Atoi(portRaw.Port())

	dsn := fmt.Sprintf(
		"user=%s dbname=%s host=%s port=%d password=%s sslmode=disable",
		dbUser, dbName, host, port, dbPassword,
	)
	conn, err := testconnectors.NewTestPgConnector(ctx, dsn, "../../../../migrations")
	if err != nil {
		log.Fatalf("pgxpool.ParseConfig: %v", err)
	}
	defer conn.Close()

	repo = numerator.NewRepository(conn)
	sqlDB = conn

	m.Run()
}

func TestRepository_AddNum(t *testing.T) {
	ctx := context.Background()

	// Очищаем таблицу перед тестами
	_, err := sqlDB.Exec(ctx, "DELETE FROM public.nums")
	require.NoError(t, err)

	tests := []struct {
		ctx         context.Context
		name        string
		inputNum    int64
		expectError error
		preAdjust   func(t *testing.T)
	}{
		{
			ctx:      ctx,
			name:     "Add single number",
			inputNum: 42,
		},
		{
			ctx:      ctx,
			name:     "Add zero",
			inputNum: 0,
		},
		{
			ctx:      ctx,
			name:     "Add negative number",
			inputNum: -7,
		},
		{
			ctx:      ctx,
			name:     "Add duplicate number (allowed)",
			inputNum: 42,
			preAdjust: func(t *testing.T) {
				_, err := sqlDB.Exec(ctx, "INSERT INTO public.nums(num) VALUES ($1)", 42)
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = sqlDB.Exec(tt.ctx, "DELETE FROM public.nums")
			require.NoError(t, err)

			if tt.preAdjust != nil {
				tt.preAdjust(t)
			}

			err = repo.AddNum(ctx, tt.inputNum)
			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
			}

			var storedNum int64
			err = sqlDB.QueryRow(tt.ctx, "SELECT num FROM public.nums WHERE num=$1", tt.inputNum).Scan(&storedNum)
			require.NoError(t, err)
			require.Equal(t, tt.inputNum, storedNum)
		})
	}
}

func TestRepository_GetSortedNums(t *testing.T) {
	ctx := context.Background()

	_, err := sqlDB.Exec(ctx, "DELETE FROM public.nums")
	require.NoError(t, err)

	tests := []struct {
		ctx         context.Context
		name        string
		insertNums  []int64
		expectError error
		expectNums  []int64
		preadjust   func(t *testing.T)
	}{
		{
			ctx:         ctx,
			name:        "Empty table returns empty slice",
			insertNums:  nil,
			expectError: nil,
			expectNums:  []int64{},
		},
		{
			ctx:         ctx,
			name:        "Returns sorted numbers in descending order",
			insertNums:  []int64{5, 2, 9, 1, 7},
			expectError: nil,
			expectNums:  []int64{1, 2, 5, 7, 9},
		},
		{
			ctx:         ctx,
			name:        "Handles single number",
			insertNums:  []int64{42},
			expectError: nil,
			expectNums:  []int64{42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = sqlDB.Exec(tt.ctx, "DELETE FROM public.nums")
			require.NoError(t, err)

			if tt.preadjust != nil {
				tt.preadjust(t)
			}

			for _, n := range tt.insertNums {
				_, err = sqlDB.Exec(tt.ctx, "INSERT INTO public.nums (num) VALUES ($1)", n)
				require.NoError(t, err)
			}

			nums, err := repo.GetSortedNums(ctx)
			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectNums, nums)
		})
	}
}
