package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"auth-regapp/helpers"
	"os"
	"github.com/joho/godotenv"
)

func Start() *mongo.Client {

	//mongoUrl := helpers.EnvFileVal("MONGOURI")
	
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error: loading .env values")
	}

	var mongoUrl string
	local := os.Getenv("MONGOURI")
	dev := os.Getenv("MONGO_DEV_CREDS")

	if dev != "" {
		mongoUrl = dev
	}else{
		mongoUrl = local
	}


	
	fmt.Println("ZZZZZZZZZZZ", mongoUrl)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl))
	if err != nil {
		log.Fatal("DB client ERROR:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()	
	  

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("DB connection ERROR:", err)
	}

	fmt.Println("Connected to MongoDB successfully!")
	return client

}

var Client *mongo.Client = Start()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("regapp").Collection(collectionName)
	return collection
}


