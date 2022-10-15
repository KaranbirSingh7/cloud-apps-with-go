package storage

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

// Database is the relational storage abstraction.
type Database struct {
	DB                    *sqlx.DB
	host                  string
	port                  int
	user                  string
	password              string
	name                  string
	maxOpenConnections    int
	maxIdleConnections    int
	connectionMaxLifetime time.Duration
	connectionMaxIdleTime time.Duration
	log                   *zap.Logger
}

// NewDatabaseOptions for NewDatabase.
type NewDatabaseOptions struct {
	Host                  string
	Port                  int
	User                  string
	Password              string
	Name                  string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
	Log                   *zap.Logger
}

func NewDatabase(opts NewDatabaseOptions) *Database {
	// no logs for DB
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}
	return &Database{
		host:                  opts.Host,
		port:                  opts.Port,
		user:                  opts.User,
		password:              opts.Password,
		name:                  opts.Name,
		maxOpenConnections:    opts.MaxOpenConnections,
		maxIdleConnections:    opts.MaxIdleConnections,
		connectionMaxLifetime: opts.ConnectionMaxLifetime,
		connectionMaxIdleTime: opts.ConnectionMaxIdleTime,
		log:                   opts.Log,
	}
}

// createDBConnectionString create a DB string (aka ConnectionURI)
func (d *Database) createDBConnectionString(withPassword bool) string {
	password := d.password

	//this ensures that we are masking password field, useful for logging/debugging
	if !withPassword {
		password = "xxx"
	}
	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", d.user, password, d.host, d.port, d.name)
}

// Connect takes care of connecting to DB and setting connections
func (d *Database) Connect() error {
	d.log.Info("connecting to database", zap.String("url", d.createDBConnectionString(false)))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	d.DB, err = sqlx.ConnectContext(ctx, "pgx", d.createDBConnectionString(true))
	if err != nil {
		return err
	}

	d.log.Debug("Setting connection pool options",
		zap.Int("max open connections", d.maxOpenConnections),
		zap.Int("max idle connections", d.maxIdleConnections),
		zap.Duration("connection max lifetime", d.connectionMaxLifetime),
		zap.Duration("connection max idle time", d.connectionMaxIdleTime),
	)

	// set timelines on DB level
	d.DB.SetMaxOpenConns(d.maxOpenConnections)
	d.DB.SetMaxIdleConns(d.maxIdleConnections)
	d.DB.SetConnMaxLifetime(d.connectionMaxLifetime)
	d.DB.SetConnMaxIdleTime(d.connectionMaxIdleTime)

	return nil
}

// Ping check if database is up and responding to queries
func (d *Database) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// check if DB is up
	if err := d.DB.PingContext(ctx); err != nil {
		return err
	}

	//check if it can handle queries
	_, err := d.DB.ExecContext(ctx, `select `)

	return err
}
