package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Client struct {
	db *mongo.Database
}

func (c *Client) Database() *mongo.Database {
	return c.db
}

var DefaultClient *Client

func Connect(url string) error {
	client, err := NewClient(url)
	if err != nil {
		return err
	}
	DefaultClient = client
	return nil
}

func NewClient(url string) (*Client, error) {
	connString, err := connstring.ParseAndValidate(url)
	if err != nil {
		return nil, err
	}

	uriOptions := options.Client().ApplyURI(connString.String())

	client, err := mongo.Connect(context.TODO(), uriOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(connString.Database)
	return &Client{
		db: db,
	}, nil
}

func (c *Client) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	c.db.Client().Disconnect(ctx)
}

/* Ctx returns a context with a timeout of 5 seconds by default. */
func (c *Client) Ctx(duration ...time.Duration) (context.Context, context.CancelFunc) {
	if len(duration) > 0 {
		return context.WithTimeout(context.Background(), duration[0])
	}
	return context.WithTimeout(context.Background(), time.Second*5)
}

func (c *Client) Collection(coll collectionName) *mongo.Collection {
	return c.db.Collection(string(coll))
}
