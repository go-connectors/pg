package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v9"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())
	return nil
}

// ErrNoRows is returned by QueryOne and ExecOne when query returned zero rows
// but at least one row is expected.
var ErrNoRows = pg.ErrNoRows

// ErrLostConnection is returned when connection to database was lost.
var ErrLostConnection = errors.New("connection to postgres database is lost")

// In accepts a slice and returns a wrapper that can be used with PostgreSQL
// IN operator:
//
//    Where("id IN (?)", pg.In([]int{1, 2, 3, 4}))
//
// produces
//
//    WHERE id IN (1, 2, 3, 4)
var In = pg.In

// DB is a database handle representing a pool of zero or more
// underlying connections. It's safe for concurrent use by multiple
// goroutines.
type DB struct {
	*pg.DB
	Config Config
}

// NewDB creates new connection to postgres using pg.v9 driver.
func NewDB(cfg *Config) (*DB, error) {
	pgdb := pg.Connect(&pg.Options{
		Addr:        cfg.Host + ":" + cfg.Port,
		User:        cfg.User,
		Password:    cfg.Password,
		Database:    cfg.Database,
		PoolSize:    cfg.PoolSize,
		PoolTimeout: time.Hour,
	})

	if cfg.Debug {
		pgdb.AddQueryHook(dbLogger{})
	}

	db := DB{DB: pgdb, Config: *cfg}
	if _, err := db.GetServerTime(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &db, nil
}

// IsConnected checks connection status to database.
func (db *DB) IsConnected() bool {
	if db == nil {
		return false
	}

	if _, err := db.GetServerTime(); err != nil {
		return false
	}

	return true
}

// GetServerTime gets and returns database server time.
func (db *DB) GetServerTime() (time.Time, error) {
	var st time.Time

	if db.DB == nil {
		return st, ErrLostConnection
	}

	_, err := db.QueryOne(pg.Scan(&st), "SELECT now()")

	return st, err
}

// Close closes connection to database.
func (db *DB) Close() error {
	if db.DB == nil {
		return nil
	}

	return db.DB.Close()
}
