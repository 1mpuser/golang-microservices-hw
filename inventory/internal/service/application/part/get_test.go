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
	"github.com/1mpuser/inventory/internal/repository/record"
	partService "github.com/1mpuser/inventory/internal/service/application/part"
	"github.com/1mpuser/inventory/internal/service/application/part/mocks"
)

func TestGet(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()

		partUUID = uuid.New()

		errRepo = errors.New("ошибка хранилища")

		part = record.PartRecord{
			UUID:          partUUID.String(),
			Name:          "Алюминиевый корпус",
			Price:         500000,
			PartType:      "HULL",
			StockQuantity: 10,
			Properties:    []byte("{}"),
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
			wantErr:   errs.ErrInvalidUUID,
		},
		{
			name: "деталь не найдена",
			uuid: partUUID.String(),
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					Get(mock.Anything, partUUID).
					Return(record.PartRecord{}, errs.ErrPartNotFound)
			},
			wantErr: errs.ErrPartNotFound,
		},
		{
			name: "ошибка репозитория",
			uuid: partUUID.String(),
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					Get(mock.Anything, partUUID).
					Return(record.PartRecord{}, errRepo)
			},
			wantErr: errRepo,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			partRepo := mocks.NewPartRepository(t)

			tc.setupMock(partRepo)

			svc := partService.NewService(partRepo, nil, nil)
			result, err := svc.Get(ctx, tc.uuid)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Empty(t, result.UUID())
			} else {
				require.NoError(t, err)
				assert.Equal(t, partUUID.String(), result.UUID())
			}
		})
	}
}
