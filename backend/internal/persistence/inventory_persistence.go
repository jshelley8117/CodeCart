package persistence

import "database/sql"

type InventoryPersistence struct {
	DbHandle *sql.DB
}

func NewInventoryPersistence(dbHandle *sql.DB) InventoryPersistence {
	return InventoryPersistence{
		DbHandle: dbHandle,
	}
}
