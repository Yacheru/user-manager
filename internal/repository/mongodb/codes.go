package mongodb

import (
	"context"
	"errors"
	"user-manager/internal/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"user-manager/init/logger"
	"user-manager/pkg/constants"
)

type Codes struct {
	coll *mongo.Collection
}

func NewCodes(coll *mongo.Collection) *Codes {
	return &Codes{coll: coll}
}

func (c *Codes) NewCode(ctx context.Context, code string, reward int) error {
	doc := bson.M{"code": code, "reward": reward}

	if _, err := c.coll.InsertOne(ctx, doc); err != nil {
		logger.Error(err.Error(), constants.MongoCategory)
		return err
	}

	return nil
}

func (c *Codes) FindCode(ctx context.Context, code string) (*entities.Code, error) {
	var codeEntity = new(entities.Code)
	filter := bson.M{"code": code}

	if err := c.coll.FindOne(ctx, filter).Decode(codeEntity); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, constants.CodeNotFoundError
		}
		logger.Error(err.Error(), constants.MongoCategory)
		return nil, err
	}

	return codeEntity, nil
}

func (c *Codes) RemoveCode(ctx context.Context, code string) error {
	filter := bson.M{"code": code}

	res, err := c.coll.DeleteOne(ctx, filter)
	if err != nil {
		logger.Error(err.Error(), constants.MongoCategory)
		return err
	}

	if res.DeletedCount == 0 {
		return constants.CodeNotFoundError
	}

	return nil
}
