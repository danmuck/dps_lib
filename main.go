package main

// import (
// 	"os"

// 	"github.com/danmuck/dps_lib/logs"
// 	"github.com/danmuck/dps_lib/mongo_client"
// 	"github.com/danmuck/dps_lib/server"

// 	"go.mongodb.org/mongo-driver/mongo"

// 	"github.com/joho/godotenv"
// )

// func init() {
// 	// tries to load .env, but won’t crash if it’s missing
// 	if err := godotenv.Load(); err != nil {
// 		logs.Info("Warning: no .env file found, relying on environment variables")
// 	}
// }

// // var CLIENT *mongo.MongoClient
// var VERSION = os.Getenv("VERSION")
// var MONGO_URI = os.Getenv("MONGO_URI") // e.g., "mongodb://localhost:27017"
// var MONGO_DB = os.Getenv("MONGO_DB")
// var MONGO_USER = os.Getenv("MONGO_USER")
// var MONGO_PASSWORD = os.Getenv("MONGO_PASSWORD")
// var USERS_DB *mongo.Collection

// func init() {
// 	var err error
// 	client, err := mongo_client.NewMongoStore(MONGO_URI, MONGO_DB)
// 	if err != nil {
// 		logs.Fatal("failed to connect to MongoDB: %v", err)
// 	}
// 	users := "users_" + VERSION // e.g., "usersv1"
// 	USERS_DB = client.Collection(users)
// 	logs.Info("Connected to MongoDB, collection: %s", USERS_DB.Name())
// }

// func main() {
// 	logs.Dev("test")
// 	s := server.NewHTTPServer()
// 	logs.Dev("Starting server with configuration: %+v", s)
// 	err := s.Start()
// 	if err != nil {
// 		logs.Err("Error starting server: %v", err)
// 		return
// 	}
// }
