package database

import (
	"context"
	"time"

	"github.com/cufee/aftermath-core/internal/core/utils"
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

var DefaultClient *Client = NewClient()

func NewClient() *Client {
	connString, err := connstring.ParseAndValidate(utils.MustGetEnv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	uriOptions := options.Client().ApplyURI(connString.String())

	client, err := mongo.Connect(context.TODO(), uriOptions)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	db := client.Database(connString.Database)
	return &Client{
		db: db,
	}
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
