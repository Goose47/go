package tasks

import (
	"Goose47/storage/config"
	"Goose47/storage/db"
	"Goose47/storage/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"path"
	"time"
)

func Init() {
	go checkExpiredItems()
}

func checkExpiredItems() {
	for {
		time.Sleep(time.Second * 1)

		cur, err := db.Client.Database("storage").Collection("storage").
			Find(
				context.TODO(),
				bson.D{
					{"exp", bson.D{{"$lte", time.Now().Unix()}}},
					{"exp", bson.D{{"$gt", 0}}},
				},
				options.Find().SetProjection(bson.D{{"_id", 1}, {"path", 1}}),
			)

		if err != nil {
			log.Println(err)
			continue
		}

		var results []models.StorageItem
		if err = cur.All(context.TODO(), &results); err != nil {
			log.Println(err)
			continue
		}

		for _, result := range results {
			nextPath := path.Join(config.FSConfig.Base, result.Path)
			if _, err := os.Stat(nextPath); os.IsNotExist(err) {
				continue
			}
			if err := os.Remove(nextPath); err != nil {
				log.Println(err)
				continue
			}
		}

		ids := make([]string, len(results))
		for i, r := range results {
			ids[i] = r.Key
		}

		_, err = db.Client.Database("storage").Collection("storage").
			DeleteMany(
				context.TODO(),
				bson.D{{"_id", bson.D{{"$in", ids}}}},
			)

		if err != nil {
			log.Println(err)
			continue
		}

		//todo move db.Client.Database("storage").Collection("storage") to config and function
		//todo what is context.TODO()???
		//todo OPENAPI SWAGGER
		//todo README
	}
}
