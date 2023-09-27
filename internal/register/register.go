package register

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

const dateFormat = `2006-01-02`

type Register struct {
	log *logrus.Logger
	db  *sql.DB
	ctx context.Context
}

func New(log *logrus.Logger, dbFN string) *Register {
	var reg Register

	reg.log = log
	reg.ctx = context.Background()

	db, err := sql.Open("sqlite", dbFN)
	if err != nil {
		log.Fatalf("database.new: Cannot establish connection to %s: %v", dbFN, err)
	}
	reg.db = db
	return &reg
}

func (reg *Register) Close() {
	reg.db.Close()
}

// Execute database query
func (reg *Register) executeTx(query string, params []any) (int64, error) {
	tx, err := reg.db.Begin()
	if err != nil {
		return -1, err
	}

	needsRollback := true
	defer cleanupTransaction(tx, &needsRollback)

	result, err := tx.ExecContext(reg.ctx, query, params...)
	if err != nil {
		return -1, err
	}

	needsRollback = false
	return result.LastInsertId()
}

func cleanupTransaction(tx *sql.Tx, needsRollback *bool) error {
	if *needsRollback {
		return tx.Rollback()
	}
	return tx.Commit()
}
