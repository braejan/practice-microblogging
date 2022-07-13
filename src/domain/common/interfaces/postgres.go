package interfaces

import (
	"context"
	"database/sql"
)

type IPostgres interface {
	Open() (DB *sql.DB, err error)
	BeginTx(ctx context.Context) (tx *sql.Tx, err error)
	CommitTx(tx *sql.Tx) (err error)
	RollbackTx(tx *sql.Tx) (err error)
}
