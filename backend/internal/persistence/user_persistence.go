package persistence

import (
	"context"

	"github.com/jshelley8117/CodeCart/internal/model"
)


type UserPersistence struct {
}

func NewUserPersistence() UserPersistence {
	return UserPersistence{}
}

func (up UserPersistence) PersistCreateUser(ctx context.Context, userDomain model.User, pgHandle ) error {

}