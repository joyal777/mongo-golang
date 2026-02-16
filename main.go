package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joyal777/mongo-golang/controllers"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := httprouter.New()

	// Get the modern MongoDB Client
	client := getClient()

	uc := controllers.NewUserController(client)

	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	r.PUT("/user/:id", uc.UpdateUser)

	fmt.Println("Server started at localhost:9000")
	log.Fatal(http.ListenAndServe("localhost:9000", r))
}

func getClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect using the official driver
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Ping to confirm connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	fmt.Println("Successfully connected to MongoDB 8.0!")
	return client
}
