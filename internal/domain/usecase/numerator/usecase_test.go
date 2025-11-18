package numerator

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUsecase_AddNum(t *testing.T) {
	ctx := context.Background()

	serviceErr := errors.New("service error")

	type args struct {
		ctx context.Context
		num int64
	}
	type mockSetup struct {
		repository func(repository *MockRepository)
	}

	tests := []struct {
		name      string
		args      args
		mockSetup mockSetup
		want      []int64
		wantErr   error
	}{
		{
			name: "successful add and get sorted nums",
			args: args{
				ctx: ctx,
				num: 5,
			},
			mockSetup: mockSetup{
				repository: func(repo *MockRepository) {
					repo.EXPECT().AddNum(mock.Anything, int64(5)).Return(nil)
					repo.EXPECT().GetSortedNums(mock.Anything).Return([]int64{9, 5, 2}, nil)
				},
			},
			want:    []int64{9, 5, 2},
			wantErr: nil,
		},
		{
			name: "AddNum returns error",
			args: args{
				ctx: ctx,
				num: 5,
			},
			mockSetup: mockSetup{
				repository: func(repo *MockRepository) {
					repo.EXPECT().AddNum(mock.Anything, int64(5)).Return(serviceErr)
				},
			},
			want:    nil,
			wantErr: serviceErr,
		},
		{
			name: "GetSortedNums returns error",
			args: args{
				ctx: ctx,
				num: 5,
			},
			mockSetup: mockSetup{
				repository: func(repo *MockRepository) {
					repo.EXPECT().AddNum(mock.Anything, int64(5)).Return(nil)
					repo.EXPECT().GetSortedNums(mock.Anything).Return(nil, serviceErr)
				},
			},
			want:    nil,
			wantErr: serviceErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)

			if tt.mockSetup.repository != nil {
				tt.mockSetup.repository(mockRepo)
			}

			u := &Usecase{
				numeratorRepository: mockRepo,
			}

			got, err := u.AddNum(tt.args.ctx, tt.args.num)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
