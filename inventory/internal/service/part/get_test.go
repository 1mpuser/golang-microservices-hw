package part_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/1mpuser/inventory/internal/errors"
	"github.com/1mpuser/inventory/internal/model"
	"github.com/1mpuser/inventory/internal/repository/record"
	partService "github.com/1mpuser/inventory/internal/service/part"
	"github.com/1mpuser/inventory/internal/service/part/mocks"
)

func TestGet(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()

		partUUID = uuid.New()

		errRepo = errors.New("ошибка хранилища")

		part = record.Part{
			UUID:          partUUID,
			Name:          "Алюминиевый корпус",
			Price:         500000,
			PartType:      model.PartTypeHull,
			StockQuantity: 10,
		}
	)

	tests := []struct {
		name      string
		uuid      string
		setupMock func(repo *mocks.PartRepository)
		wantErr   error
	}{
		{
			name: "деталь найдена",
			uuid: partUUID.String(),
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					Get(mock.Anything, partUUID).
					Return(part, nil)
			},
			wantErr: nil,
		},
		{
			name:      "неверный формат uuid",
			uuid:      "не-uuid",
			setupMock: func(_ *mocks.PartRepository) {},
			wantErr:   errs.ErrInvalidFormat,
		},
		{
			name: "деталь не найдена",
			uuid: partUUID.String(),
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					Get(mock.Anything, partUUID).
					Return(record.Part{}, errs.ErrPartNotFound)
			},
			wantErr: errs.ErrPartNotFound,
		},
		{
			name: "ошибка репозитория",
			uuid: partUUID.String(),
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					Get(mock.Anything, partUUID).
					Return(record.Part{}, errRepo)
			},
			wantErr: errRepo,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			partRepo := mocks.NewPartRepository(t)

			tc.setupMock(partRepo)

			svc := partService.NewService(partRepo)
			result, err := svc.Get(ctx, tc.uuid)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Empty(t, result.UUID)
			} else {
				require.NoError(t, err)
				assert.Equal(t, partUUID.String(), result.UUID)
			}
		})
	}
}
