package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/froggy-12/mooshroombase/api"
	"github.com/froggy-12/mooshroombase/config"
	"github.com/froggy-12/mooshroombase/docker"
	"github.com/froggy-12/mooshroombase/utils"
	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient   *mongo.Client
	mariaDBClient *sql.DB
	redisClient   *redis.Client
)

func main() {
	utils.DebugLogger("main", "initializing Configuration Files 📁📁📂")
	config.Configs = config.InitConfigs()
	utils.DebugLogger("main", "Done 🍃👍")

	config.CheckIfFieldsAreEmpty(config.Configs)
	if config.Configs.Authentication {
		if config.Configs.GithubKey == "" || config.Configs.GithubSecret == "" || config.Configs.GoogleKey == "" || config.Configs.GoogleSecret == "" {
			log.Fatal("Invalid Configurations Authentication is enabled no configurations handled")
		}
	}

	// initializing docker configs
	utils.DebugLogger("main", "Starting Docker Configurations")
	docker.Init()

	// database connections and api servers
	utils.DebugLogger("main", "Starting Database Connections")

	for _, database := range config.Configs.RunningDatabaseContainers {
		switch database {
		case "mongodb":
			mongoURI := fmt.Sprintf("mongodb://%v:%v@localhost:27017", config.Configs.MongoDBUsername, config.Configs.MongoDBPassword)
			client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
			if err != nil {
				log.Fatal("Failed to connect with MongoDB ", err.Error())
			}
			utils.DebugLogger("main", "Connection with MongoDB 🍃🍃😊 successfull..")
			mongoClient = client
		case "redis":
			options := &redis.Options{
				Addr: fmt.Sprintf("localhost:%v", "6379"),
			}
			client := redis.NewClient(options)
			_, err := client.Ping(context.Background()).Result()
			if err != nil {
				log.Fatal("Error Connection with Redis: ", err.Error())
			}
			utils.DebugLogger("main", "Connection with Redis 🔴 successfull..")
			redisClient = client

		case "mariadb":
			cfg := mysql.Config{
				User:                 "root",
				Passwd:               config.Configs.MariaDBRootPassword,
				Addr:                 "localhost:" + "3306",
				Net:                  "tcp",
				AllowNativePasswords: true,
				ParseTime:            true,
			}
			db, err := sql.Open("mysql", cfg.FormatDSN())
			if err != nil {
				log.Fatal("Error Connecting to MariaDB: ", err.Error())
			}
			err = db.Ping()
			if err != nil {
				log.Fatal("Error Connecting to MariaDB: ", err.Error())
			}
			utils.DebugLogger("main", "MariaDB 🐬 Connection Successful 👍")
			mariaDBClient = db
		}
	}

	if config.Configs.Authentication {
		if config.Configs.PrimaryDB == "mongodb" {
			utils.DebugLogger("main", "detected mongodb as primary db indexing some models")
			database := mongoClient.Database("mooshroombase")
			usersCollection := database.Collection("users")
			_, err := usersCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
				Keys:    bson.M{"email": 1},
				Options: options.Index().SetUnique(true),
			})

			if err != nil {
				log.Fatal(err)
			}

			_, err = usersCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
				Keys:    bson.M{"username": 1},
				Options: options.Index().SetUnique(true),
			})
			if err != nil {
				log.Fatal(err)
			}

		}
	}

	utils.DebugLogger("main", "Database Connections are successfull 😊😊")

	utils.DebugLogger("main", "Starting THE API SERVER!!!!")

	server := api.APIServer(":6644", mongoClient, redisClient, mariaDBClient)
	err := server.Start()

	if err != nil {
		log.Fatal("Error Starting the API Server")
	}
}
