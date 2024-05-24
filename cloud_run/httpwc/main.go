package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
)

type TimeControl struct {
	GameLengthInSeconds int `json:"game_length_second"`
	IncrementInSeconds  int `json:"increment_second"`
}

func main() {

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("Error initializing firebase app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error initializing firestore client: %v", err)
	}
	defer client.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", getCollection(ctx, client, "time-control"))
	http.HandleFunc("/collection", handleCollection(ctx, client, "time-control"))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Could not listen")
	}
}

func handleCollection(ctx context.Context, client *firestore.Client, collection string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addCollection(ctx, client, collection, w, r)
		} else if r.Method == http.MethodPatch {
			updateCollection(ctx, client, collection, w, r)
		}
	}
}

func updateCollection(ctx context.Context, client *firestore.Client, collection string, w http.ResponseWriter, r *http.Request) {
	docID := r.URL.Query().Get("id")
	if docID == "" {
		http.Error(w, "Missing document ID in request URL", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := client.Collection(collection).Doc(docID).Set(ctx, updates, firestore.MergeAll)
	if err != nil {
		http.Error(w, "Error updating collection", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Updated collection with new Settings")
}

func addCollection(ctx context.Context, client *firestore.Client, collection string, w http.ResponseWriter, r *http.Request) {
	var doc TimeControl
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, _, err := client.Collection(collection).Add(ctx, map[string]interface{}{
		"game_length_second": doc.GameLengthInSeconds,
		"increment_second":   doc.IncrementInSeconds,
	})
	if err != nil {
		http.Error(w, "Error adding time control document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Created time control document successfully")
}

func getCollection(ctx context.Context, client *firestore.Client, collection string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		iter := client.Collection(collection).Documents(ctx)
		defer iter.Stop()

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				http.Error(w, "Error fetching Document", http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "Document data: %v/n", doc.Data())
		}
	}
}
