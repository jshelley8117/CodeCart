package persistence

import (
	"context"
	"database/sql"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type UserPersistence struct {
	DbHandle *sql.DB
}

func NewUserPersistence(dbHandle *sql.DB) UserPersistence {
	return UserPersistence{
		DbHandle: dbHandle,
	}
}

func (up UserPersistence) PersistCreateUser(ctx context.Context, userDomain model.User) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistCreateUser")
	query := `
		INSERT INTO users (email, created_at, updated_at, is_active, customer_id, auth_id, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := up.DbHandle.ExecContext(
		ctx,
		query,
		userDomain.Email,
		userDomain.CreatedAt,
		userDomain.UpdatedAt,
		userDomain.IsActive,
		userDomain.CustomerId,
		userDomain.AuthId,
		userDomain.Role,
	)
	if err != nil {
		z.Error("ExecContext failed: %w", zap.Error(err))
		return err
	}

	return nil
}
