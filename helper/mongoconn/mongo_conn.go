package mongoconn

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// Connector - MongoDB 连接 适配器
type Connector struct {
	client   *mongo.Client
	database *mongo.Database
	timeout  time.Duration // 此处 timeout 时间长度尤其要 ⚠️
}

// finalizer is the destructor for adapter.
/*
func finalizer(c *Connector) {
	c.close()
}
*/

// NewConn - 连接 适配器 工厂
func NewConn(url string, timeout time.Duration) (*Connector, error) {
	log.Println("url:", url)

	// Parse and validate url before apply it.
	connString, err := connstring.ParseAndValidate(url)

	if err != nil {
		return nil, err
	}

	clientOption := options.Client().ApplyURI(url)

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, clientOption)

	if err != nil {
		return nil, err
	}

	// Get database name from connString.
	if connString.Database == "" {
		return nil, errors.New("Database name not provided")
	}

	c := &Connector{client, client.Database(connString.Database), timeout * time.Second}

	// runtime.SetFinalizer(c, finalizer)

	return c, nil
}

// Close - Disconnect
func (c *Connector) Close() error {
	log.Println("准备与数据库断开连接")
	ctx, cancel := context.WithTimeout(context.TODO(), c.timeout)
	defer cancel()
	return c.client.Disconnect(ctx)
}

// Close - 关闭连接
/*
func (c *Connector) close() {
	log.Println("准备与数据库断开连接")
	ctx, cancel := context.WithTimeout(context.TODO(), c.timeout)
	defer cancel()
	err := c.client.Disconnect(ctx)
	if err != nil {
		log.Println("与数据库断开连接失败")
	}
}
*/

// Ping - 建立与数据库的连接
func (c *Connector) Ping() error {
	err := c.client.Ping(context.TODO(), nil)

	if err != nil {
		return err
	}

	log.Println("Connected to MongoDB!")
	return nil
}

// GetColl - Get spec collection
func (c *Connector) GetColl(collname string) *mongo.Collection {
	return c.database.Collection(collname)
}
