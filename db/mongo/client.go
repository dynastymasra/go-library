package mongo

import (
	"context"
	"os"

	"github.com/matryer/resync"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	once   resync.Once
)

// Config holds the configuration settings for connecting to a MongoDB instance.
// It contains the following fields:
// - URI is the connection string for the MongoDB instance.
// - Username is the username for authenticating to the MongoDB instance.
// - Password is the password for authenticating to the MongoDB instance.
// - MaxPoolSize specifies the maximum number of connections in the connection pool.
// - MinPoolSize specifies the minimum number of connections in the connection pool.
type Config struct {
	URI, Username, Password  string
	MaxPoolSize, MinPoolSize uint64
}

// Connect initializes the MongoDB client and establishes a connection to the database.
// It uses a sync.Once to ensure the client is only initialized once.
// If the hostname cannot be determined, it defaults to "localhost".
// The method sets up the client with the provided configuration options and pings the database to verify the connection.
//
// Parameters:
// - ctx: The context to use for the connection.
//
// Returns:
// - error: An error if the connection or ping fails, otherwise nil.
func (c Config) Connect(ctx context.Context) error {
	var err error
	once.Do(func() {
		hostname, err := os.Hostname()
		if hostname == "" || err != nil {
			hostname = "localhost"
		}

		opts := options.Client().ApplyURI(c.URI).SetAppName(hostname).SetMaxPoolSize(c.MaxPoolSize).
			SetMinPoolSize(c.MinPoolSize).SetAuth(options.Credential{
			Username: c.Username,
			Password: c.Password,
		})
		if err = opts.Validate(); err != nil {
			return
		}

		client, err = mongo.Connect(ctx, opts)
		if err != nil {
			return
		}

		err = client.Ping(ctx, nil)
	})

	return err
}

// Client returns the MongoDB client instance.
// This method provides access to the initialized MongoDB client.
//
// Returns:
// - *mongo.Client: The MongoDB client instance.
func (c Config) Client() *mongo.Client {
	return client
}

// Ping checks the connection to the MongoDB instance by sending a ping command.
//
// Returns:
// - error: An error if the ping fails, otherwise nil.
func (c Config) Ping() error {
	return client.Ping(context.Background(), nil)
}

// SetClient sets the MongoDB client instance to the provided connection.
//
// Parameters:
// - conn: The MongoDB client instance to set.
func (c Config) SetClient(conn *mongo.Client) {
	client = conn
}

// Reset reinitialize the MongoDB client by resetting the sync.Once instance and reconnecting to the database.
// It first resets the sync.Once instance to allow reinitialization.
// Then, it attempts to reconnect to the MongoDB instance using the Connect method.
//
// Returns:
// - error: An error if the reconnection fails, otherwise nil.
func (c Config) Reset() error {
	once.Reset()
	if err := c.Connect(context.Background()); err != nil {
		return err
	}

	return nil
}
