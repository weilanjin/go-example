package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type Transaction struct {
	opts *sql.TxOptions
	fun  func(ctx context.Context, tx *sql.Tx) error
}

func Tx(fn func(ctx context.Context, tx *sql.Tx) error) *Transaction {
	return &Transaction{
		fun: fn,
	}
}

func (t *Transaction) WithOpts(opts *sql.TxOptions) *Transaction {
	t.opts = opts
	return t
}

func (t *Transaction) Exec(ctx context.Context, db *sql.DB) (err error) {
	tx, err := db.BeginTx(ctx, t.opts)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				slog.Error(fmt.Sprintf("tx.Rollback() err: %v", rbErr))
			}
		}
	}()

	if err = t.fun(ctx, tx); err != nil {
		return err
	}
	return tx.Commit()
}
