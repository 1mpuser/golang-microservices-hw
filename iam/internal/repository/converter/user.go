package converter

import (
	"github.com/1mpuser/iam/internal/model"
	"github.com/1mpuser/iam/internal/repository/record"
	"github.com/google/uuid"
)

func ModelToRecord(r record.User) model.User {
	uidTrasformed := uuid.MustParse(r.UUID)

	return model.User{
		UUID:         uidTrasformed,
		Login:        r.Login,
		PasswordHash: r.PasswordHash,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    &r.UpdatedAt,
	}
}

func RecordToModel(m model.User) record.User {
	return record.User{
		UUID:         m.UUID.String(),
		Login:        m.Login,
		PasswordHash: m.PasswordHash,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    *m.UpdatedAt,
	}
}
