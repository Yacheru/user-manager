package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"user-manager/init/config"
	"user-manager/init/logger"
	"user-manager/pkg/constants"
)

func InitMongoDB(ctx context.Context, cfg *config.Config) (*mongo.Collection, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.MongodbURI).SetServerAPIOptions(serverAPI)

	logger.Debug("create mongo client", constants.MongoCategory)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.Error(err.Error(), constants.MongoCategory)
		return nil, err
	}

	db := client.Database(cfg.MongodbDatabase)

	logger.DebugF("pinging mongo database (%s)", constants.MongoCategory, cfg.MongodbDatabase)

	if err := db.RunCommand(ctx, bson.D{{"ping", 1}}).Err(); err != nil {
		logger.Error(err.Error(), constants.MongoCategory)
		return nil, err
	}

	coll := db.Collection(cfg.MongodbCollection)

	logger.Info("successfully connected to MongoDB!!!", constants.MongoCategory)

	return coll, nil
}
