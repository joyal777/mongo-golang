package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/joyal777/mongo-golang/controllers"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := httprouter.New()

	// Get the modern MongoDB Client with connection pooling
	client := getClient()

	// Create controller with the client
	uc := controllers.NewUserController(client)

	// Setup routes
	r.GET("/users", uc.GetAllUsers)
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	r.PUT("/user/:id", uc.UpdateUser)

	// Health check endpoint to demonstrate concurrency
	r.GET("/health", healthCheck)

	// Demo endpoint for concurrent operations
	r.GET("/demo/concurrent", demoConcurrentOperations(client))

	fmt.Println("Server started at localhost:9000")
	log.Fatal(http.ListenAndServe("localhost:9000", r))
}

func getClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure connection pool for better concurrency
	clientOptions := options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetMinPoolSize(10).                 // Minimum connections in pool
		SetMaxPoolSize(100).                // Maximum connections in pool
		SetMaxConnIdleTime(5 * time.Minute) // Close idle connections

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping to confirm connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	fmt.Println("Successfully connected to MongoDB with connection pooling!")
	return client
}

// Simple health check - demonstrates goroutine for logging
func healthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Start a goroutine for async logging (doesn't block response)
	go func() {
		log.Printf("Health check accessed from %s at %v",
			r.RemoteAddr, time.Now().Format(time.RFC3339))
	}()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "healthy", "time": "%s"}`, time.Now().Format(time.RFC3339))
}

// Demo function showing various concurrency patterns
func demoConcurrentOperations(client *mongo.Client) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Channel to collect results
		resultCh := make(chan string, 3) // Buffered channel
		errorCh := make(chan error, 3)

		// WaitGroup to wait for all goroutines
		var wg sync.WaitGroup

		// Start 3 concurrent database operations
		operations := []string{"users", "products", "orders"}

		for _, collection := range operations {
			wg.Add(1) // Add to waitgroup before starting goroutine

			// Launch goroutine for each collection
			go func(collName string) {
				defer wg.Done() // Signal completion when done

				// Simulate database operation
				count, err := client.Database("testdb").Collection(collName).CountDocuments(ctx, map[string]interface{}{})
				if err != nil {
					errorCh <- fmt.Errorf("error counting %s: %v", collName, err)
					return
				}

				// Send result through channel
				resultCh <- fmt.Sprintf("%s: %d documents", collName, count)
			}(collection) // Pass collection name to avoid closure issues
		}

		// Start another goroutine to close channels when all work is done
		go func() {
			wg.Wait()       // Wait for all operations
			close(resultCh) // Close result channel
			close(errorCh)  // Close error channel
		}()

		// Collect results (this runs concurrently with the operations)
		var results []string
		var errors []string

		// Use select to handle multiple channels
		for {
			select {
			case result, ok := <-resultCh:
				if !ok {
					resultCh = nil // Channel closed
				} else {
					results = append(results, result)
				}
			case err, ok := <-errorCh:
				if !ok {
					errorCh = nil // Channel closed
				} else {
					errors = append(errors, err.Error())
				}
			}

			// Exit when both channels are closed
			if resultCh == nil && errorCh == nil {
				break
			}
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		if len(errors) > 0 {
			w.WriteHeader(http.StatusPartialContent)
			fmt.Fprintf(w, `{"results": %v, "errors": %v}`, results, errors)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"results": %v}`, results)
		}
	}
}
