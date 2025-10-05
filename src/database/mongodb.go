package database

import (
	"Monitoring-Opportunities/src/config"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDatabase(cfg config.Config) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(cfg.MongoURI).
		SetMaxPoolSize(300).                        // Tăng pool size cho high throughput
		SetMinPoolSize(50).                         // Giữ sẵn connections để tránh cold start
		SetMaxConnIdleTime(5 * time.Minute).        // Giảm idle time để giải phóng connection nhanh
		SetMaxConnecting(10).                       // Giới hạn số connection đang được thiết lập
		SetConnectTimeout(5 * time.Second).         // Timeout khi thiết lập connection
		SetSocketTimeout(30 * time.Second).         // Timeout cho mỗi operation
		SetServerSelectionTimeout(5 * time.Second). // Timeout chọn server
		SetRetryReads(true).                        // Tự động retry cho read operations
		SetRetryWrites(true)                        // Tự động retry cho write operations

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client.Database(cfg.DBName), nil
}
