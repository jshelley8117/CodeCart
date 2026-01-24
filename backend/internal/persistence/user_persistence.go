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
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered PersistCreateUser")
	query := `
		INSERT INTO users (email, created_at, updated_at, is_active, customer_id, gc_auth_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := up.DbHandle.ExecContext(
		ctx,
		query,
		userDomain.Email,
		userDomain.CreatedAt,
		userDomain.UpdatedAt,
		userDomain.IsActive,
		userDomain.CustomerId,
		userDomain.GCAuthId,
	)
	if err != nil {
		zLog.Error("ExecContext failed: %w", zap.Error(err))
		return err
	}

	return nil
}
