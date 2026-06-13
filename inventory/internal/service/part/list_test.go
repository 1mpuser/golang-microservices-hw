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
	partService "github.com/1mpuser/inventory/internal/service/part"
	"github.com/1mpuser/inventory/internal/service/part/mocks"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
)

func TestList(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()

		hullUUID   = uuid.New()
		engineUUID = uuid.New()

		errRepo = errors.New("ошибка хранилища")

		hullPart = record.Part{
			UUID:     hullUUID,
			Name:     "Титановый корпус",
			PartType: inventoryv1.PartType_PART_TYPE_HULL,
		}

		enginePart = record.Part{
			UUID:     engineUUID,
			Name:     "Алюминиевый корпус",
			PartType: inventoryv1.PartType_PART_TYPE_HULL,
		}
	)

	tests := []struct {
		name      string
		uuids     []string
		partType  inventoryv1.PartType
		setupMock func(repo *mocks.PartRepository)
		wantErr   error
		wantNames []string
	}{
		{
			name:     "все детали без фильтра, сортировка по имени",
			uuids:    nil,
			partType: inventoryv1.PartType_PART_TYPE_UNSPECIFIED,
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					ListAllParts(mock.Anything).
					Return([]record.Part{hullPart, enginePart}, nil)
			},
			wantErr:   nil,
			wantNames: []string{"Алюминиевый корпус", "Титановый корпус"},
		},
		{
			name:     "фильтр по типу детали",
			uuids:    nil,
			partType: inventoryv1.PartType_PART_TYPE_HULL,
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					ListPartsByPartType(mock.Anything, inventoryv1.PartType_PART_TYPE_HULL).
					Return([]record.Part{hullPart}, nil)
			},
			wantErr:   nil,
			wantNames: []string{"Титановый корпус"},
		},
		{
			name:     "фильтр по списку uuid",
			uuids:    []string{hullUUID.String()},
			partType: inventoryv1.PartType_PART_TYPE_UNSPECIFIED,
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					ListPartsByUuids(mock.Anything, []uuid.UUID{hullUUID}).
					Return([]record.Part{hullPart}, nil)
			},
			wantErr:   nil,
			wantNames: []string{"Титановый корпус"},
		},
		{
			name:      "неверный формат uuid в фильтре",
			uuids:     []string{"не-uuid"},
			partType:  inventoryv1.PartType_PART_TYPE_UNSPECIFIED,
			setupMock: func(_ *mocks.PartRepository) {},
			wantErr:   errs.ErrInvalidFormat,
		},
		{
			name:     "одна из деталей не найдена",
			uuids:    []string{hullUUID.String()},
			partType: inventoryv1.PartType_PART_TYPE_UNSPECIFIED,
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					ListPartsByUuids(mock.Anything, []uuid.UUID{hullUUID}).
					Return(nil, errs.ErrPartNotFound)
			},
			wantErr: errs.ErrPartNotFound,
		},
		{
			name:     "ошибка репозитория",
			uuids:    nil,
			partType: inventoryv1.PartType_PART_TYPE_UNSPECIFIED,
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					ListAllParts(mock.Anything).
					Return(nil, errRepo)
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
			result, err := svc.List(ctx, tc.uuids, tc.partType)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)

			names := make([]string, 0, len(result))
			for _, p := range result {
				names = append(names, p.Name)
			}
			assert.Equal(t, tc.wantNames, names)
		})
	}
}
