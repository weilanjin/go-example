package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log/slog"
	"net/url"
)

// New 初始化mongo
func New(uri string) (*mongo.Client, error) {
	ctx := context.Background()

	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, err
	}

	// mgo日志输出：debug sql
	mgoLogOpt := event.CommandMonitor{
		Started: func(ctx context.Context, e *event.CommandStartedEvent) {
			slog.Debug("mongo start", slog.Int64("reqId", e.RequestID), slog.String("db", e.DatabaseName),
				slog.String("cmd", e.CommandName), slog.Any("data", e.Command.String()))
		},
		Succeeded: func(ctx context.Context, e *event.CommandSucceededEvent) {
			slog.Debug("mongo succeed", slog.Int64("reqId", e.RequestID), slog.String("cmd", e.CommandName), slog.Duration("time", e.Duration))
		},
		Failed: func(ctx context.Context, e *event.CommandFailedEvent) {
			slog.Info("mongo failed", slog.Int64("reqId", e.RequestID), slog.String("cmd", e.CommandName),
				slog.Duration("time", e.Duration), slog.String("reason", e.Failure))
		},
	}

	bson.DefaultRegistry.RegisterTypeEncoder(tDecimal, Decimal{})
	bson.DefaultRegistry.RegisterTypeDecoder(tDecimal, Decimal{})

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetMonitor(&mgoLogOpt))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return client, nil
}
