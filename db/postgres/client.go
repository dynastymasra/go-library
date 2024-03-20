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

// Client is a method on the Config struct that establishes a database connection using the configuration in the Config struct.
// It first constructs the data source name (dsn) using the configuration fields.
// Then it sets the log mode based on the LogMode field in the Config struct.
// It then attempts to open a connection to the database using the constructed dsn and log mode.
// If an error occurs during the connection process, it returns the error.
// If the DebugEnabled field in the Config struct is set to true, it enables debug mode on the database connection.
// It then retrieves the underlying sql.DB object from the gorm.DB object and sets the maximum idle and open connections.
// Finally, it returns the established database connection and any error that occurred during the process.
//
// Returns:
// - *gorm.DB: The established database connection.
// - error: Any error that occurred during the process.
func (c Config) Client() (*gorm.DB, error) {
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
	return db, err
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
// It first resets the once.Do function to allow the Client method to be called again.
// Then it calls the Client method to establish a new database connection.
// If an error occurs during this process, it returns the error.
// If the process is successful, it returns nil.
func (c Config) Reset() error {
	var err error

	once.Reset()
	db, err = c.Client()
	if err != nil {
		return err
	}

	return nil
}
