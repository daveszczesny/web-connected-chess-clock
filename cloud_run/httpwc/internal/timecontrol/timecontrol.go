package timecontrol

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TimeControl struct {
	GameLengthInMilliseconds int `json:"game_length_ms"`
	IncrementInMilliseconds  int `json:"increment_ms"`
}

/*
Handler function to handle collection API calls
*/
func HandleCollection(ctx context.Context, client *firestore.Client, collection string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch method := r.Method; method {
		case http.MethodGet:
			getCollection(ctx, client, collection, w)
		case http.MethodPost:
			addCollection(ctx, client, collection, w, r)
		case http.MethodPatch:
			updateCollection(ctx, client, collection, w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getCollection(ctx context.Context, client *firestore.Client, collection string, w http.ResponseWriter) {
	iter := client.Collection(collection).Documents(ctx)
	defer iter.Stop()

	var documents []map[string]interface{}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Error fetching Document", http.StatusInternalServerError)
		}
		documents = append(documents, doc.Data())
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(documents)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func updateCollection(ctx context.Context, client *firestore.Client, collection string, w http.ResponseWriter, r *http.Request) {
	docID := r.URL.Query().Get("id")
	if docID == "" {
		http.Error(w, "Missing document ID in request URL", http.StatusBadRequest)
		return
	}

	// Decode payload
	var updates map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	docRef := client.Collection(collection).Doc(docID)
	_, err = docRef.Get(ctx)
	// Check if document exists
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error checking document existence", http.StatusInternalServerError)
		return
	}

	// Update document
	_, err = docRef.Set(ctx, updates, firestore.MergeAll)
	if err != nil {
		http.Error(w, "Error updating collection", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Updated collection with new Settings")
}

func addCollection(ctx context.Context, client *firestore.Client, collection string, w http.ResponseWriter, r *http.Request) {
	var doc TimeControl
	err := json.NewDecoder(r.Body).Decode(&doc)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, _, err = client.Collection(collection).Add(ctx, map[string]interface{}{
		"game_length_ms": doc.GameLengthInMilliseconds,
		"increment_ms":   doc.IncrementInMilliseconds,
	})
	if err != nil {
		http.Error(w, "Error adding time control document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Created time control document successfully")
}
