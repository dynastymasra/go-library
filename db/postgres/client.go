package postgres

import (
	"fmt"
	"strings"

	"github.com/matryer/resync"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once resync.Once
)

// Config is a struct that holds the configuration for the database connection.
// It contains the following fields:
// - Database: the name of the database
// - Host: the host of the database
// - Username: the username to connect to the database
// - Password: the password to connect to the database
// - Params: additional parameters for the database connection
// - LogMode: the log mode for the database connection
// - Port: the port number of the database
// - MaxIdleConn: the maximum number of idle connections for the database
// - MaxOpenConn: the maximum number of open connections for the database
// - DebugEnabled: a boolean indicating whether debug mode is enabled or not
type Config struct {
	Database, Host, Username, Password, Params, LogMode string
	Port, MaxIdleConn, MaxOpenConn                      int
	DebugEnabled                                        bool
}

// Connect is a method on the Config struct that establishes a connection to the database.
// It constructs the Data Source Name (DSN) using the configuration fields and opens a connection to the database.
// The connection is established only once using singleton mechanism to ensure that the connection is not re-established multiple times.
// If the connection is successfully established, it configures the connection pool settings and enables debug mode if specified.
//
// Returns:
// - error: Any error that occurred during the connection process.
func (c Config) Connect() error {
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d %s",
			c.Username, c.Password, c.Database, c.Host, c.Port, c.Params)

		logMode := func() logger.LogLevel {
			switch strings.ToLower(c.LogMode) {
			case "error":
				return logger.Error
			case "warn":
				return logger.Warn
			case "info":
				return logger.Info
			default:
				return logger.Silent
			}
		}

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logMode()),
		})
		if err != nil {
			return
		}

		if c.DebugEnabled {
			db = db.Debug()
		}

		sqlDB, err := db.DB()
		if err != nil {
			return
		}

		sqlDB.SetMaxIdleConns(c.MaxIdleConn)
		sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	})
	return err
}

// Ping is a method on the Config struct that checks the database connection by sending a ping.
// If the database is reachable and responds to the ping, it returns nil.
// If the database is not reachable or does not respond to the ping, it returns an error.
func (c Config) Ping() error {
	conn, err := db.DB()
	if err != nil {
		return err
	}

	return conn.Ping()
}

// Close is a method on the Config struct that closes the database connection.
// It first retrieves the underlying sql.DB object from the gorm.DB object.
// If an error occurs during this process, it returns the error.
// If the retrieval is successful, it calls the Close method on the sql.DB object to close the database connection.
// If an error occurs while closing the database connection, it returns the error.
// If the database connection is successfully closed, it returns nil.
func (c Config) Close() error {
	conn, err := db.DB()
	if err != nil {
		return err
	}
	return conn.Close()
}

// SetDB is a method on the Config struct that sets the database connection.
// It takes one parameter:
// - conn: The *gorm.DB object representing the database connection to be set.
// This method does not return any value.
func (c Config) SetDB(conn *gorm.DB) {
	db = conn
}

// DB is a method on the Config struct that retrieves the current database connection.
// It does not take any parameters.
// It returns the *gorm.DB object representing the current database connection.
func (c Config) DB() *gorm.DB {
	return db
}

// Reset is a method on the Config struct that resets the database connection.
// It first resets the instance to allow the Connect method to be called again.
// Then it calls the Connect method to re-establish the database connection.
// If an error occurs during the connection process, it returns the error.
// If the connection is successfully re-established, it returns nil.
//
// Returns:
// - error: Any error that occurred during the process.
func (c Config) Reset() error {
	once.Reset()
	return c.Connect()
}
