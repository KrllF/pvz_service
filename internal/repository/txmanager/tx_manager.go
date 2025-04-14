package txmanager

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DB интерфейс функций для работы с бд
type (
	DB interface {
		Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
		QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
		Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
		GetPool() *pgxpool.Pool
	}

	// txManagerKey ключ для менеджера транзакций
	txManagerKey struct{}

	// TxManager менеджер транзакций
	TxManager struct {
		DB
	}
)

// NewTxManager конструктур менеджера транзакций
func NewTxManager(db DB) (*TxManager, error) {
	return &TxManager{DB: db}, nil
}

// RunSerializable уровень изоляции Serializable
func (m *TxManager) RunSerializable(ctx context.Context, fn func(ctxTx context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}

	return m.beginFunc(ctx, opts, fn)
}

// RunReadCommitted уровень изоляции Committed
func (m *TxManager) RunReadCommitted(ctx context.Context, fn func(ctxTx context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	}

	return m.beginFunc(ctx, opts, fn)
}

func (m *TxManager) beginFunc(ctx context.Context, opts pgx.TxOptions, fn func(ctxTx context.Context) error) error {
	tx, err := m.DB.GetPool().BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	ctx = context.WithValue(ctx, txManagerKey{}, tx)
	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetDB возвращает интерфейс DB
func (m *TxManager) GetDB(ctx context.Context) DB {
	v, ok := ctx.Value(txManagerKey{}).(DB)
	if ok && v != nil {
		return v
	}

	return m.DB
}
