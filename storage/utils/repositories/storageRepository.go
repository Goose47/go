package repositories

import (
	"Goose47/storage/db"
	"Goose47/storage/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type StorageRepository struct{}

func (*StorageRepository) FindByKey(key string) (*models.StorageItem, error) {
	var result bson.M
	err := db.Client.Database("storage").
		Collection("storage").
		FindOne(context.TODO(), bson.D{{"_id", key}}).
		Decode(&result)

	if err != nil {
		return nil, err
	}

	return createItem(result), nil
}

func (*StorageRepository) Set(key string, item *models.StorageItem) (string, error) {
	_, err := db.Client.Database("storage").
		Collection("storage").
		DeleteOne(context.TODO(), bson.D{{"_id", key}})

	if err != nil {
		return "", err
	}

	result, err := db.Client.Database("storage").
		Collection("storage").
		InsertOne(context.TODO(), bson.D{
			{"_id", key},
			{"path", item.Path},
			{"ttl", item.Ttl},
			{"originalName", item.OriginalName},
		})

	if err != nil {
		return "", err
	}

	return result.InsertedID.(string), nil
}

func (*StorageRepository) DeleteByKey(key string) (*models.StorageItem, error) {
	var result bson.M
	err := db.Client.Database("storage").
		Collection("storage").
		FindOneAndDelete(context.TODO(), bson.D{{"_id", key}}).
		Decode(&result)

	if err != nil {
		return nil, err
	}

	return createItem(result), nil
}

func createItem(result map[string]interface{}) *models.StorageItem {
	item := &models.StorageItem{}

	item.Key = result["_id"].(string)
	item.Path = result["path"].(string)
	item.Ttl = int(result["ttl"].(int32))
	item.OriginalName = result["originalName"].(string)

	return item
}
