//nolint:all
package httptest

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/KrllF/pvz_service/internal/config"
	"github.com/KrllF/pvz_service/internal/db"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

// DB интерфейс запросов
type DB interface {
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Close()
}

// TDB структура для работы с бд при тесте
type TDB struct {
	DB
}

// NewFromEnv конструктор TDB
func NewFromEnv() *TDB {
	conf, err := config.New("config.yaml")
	if err != nil {
		panic(err)
	}
	dbL, err := db.NewDB(context.Background(), conf)
	if err != nil {
		panic(err)
	}

	return &TDB{DB: dbL}
}

// SetUp добавление и очищение таблицы перед тестом
func (d *TDB) SetUp(t *testing.T, tableName ...string) {
	t.Helper()
	d.truncateTable(context.Background(), tableName...)
	d.addUser(context.Background())
	d.addOrder(context.Background())
}

// TearDown после теста
func (d *TDB) TearDown(t *testing.T) {
	t.Helper()
}

// Close закрыть пулл TDB
func (d *TDB) Close() {
	d.DB.Close()
}

// addUser добавить пользователя перед тестом
func (d *TDB) addUser(ctx context.Context) {
	_, err := d.DB.Exec(ctx, "INSERT INTO users(user_id) VALUES($1)", 1)
	if err != nil {
		panic(err)
	}
}

// addOrder добавить заказ перед тестом
func (d *TDB) addOrder(ctx context.Context) {
	queryAdd := `INSERT INTO orders(order_id,user_id,
	order_weight, order_price,
	pack_id, extra_pack_id,
	status_id, shelf_life,
	two_days_of_life, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8, $9, $10, $11);`

	twoDaysOfLife := pgtype.Timestamptz{
		Time:   time.Time{},
		Status: pgtype.Null,
	}

	_, err := d.DB.Exec(ctx, queryAdd, 1, 1, 100, 100,
		1, 1, 1, "2027-10-10T15:15:15Z", twoDaysOfLife, time.Now(), time.Now())
	if err != nil {
		panic(err)
	}
	_, err = d.DB.Exec(ctx, queryAdd, 100, 1, 100, 100,
		1, 1, 3, "2023-10-10T15:15:15Z", twoDaysOfLife, time.Now(), time.Now())
	if err != nil {
		panic(err)
	}
}

// truncateTable очистить таблицу
func (d *TDB) truncateTable(ctx context.Context, tableName ...string) {
	q := "TRUNCATE " + strings.Join(tableName, ",")
	if _, err := d.DB.Exec(ctx, q); err != nil {
		panic(err)
	}
}
