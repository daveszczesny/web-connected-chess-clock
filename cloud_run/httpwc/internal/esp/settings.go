package esp

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func HandleCollection(ctx context.Context, client *firestore.Client, collection string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch method := r.Method; method {
		case http.MethodGet:
			getCollection(ctx, client, collection, w)
		case http.MethodPatch:
			updateCollection(ctx, client, collection, w, r)
		}
	}
}

/*
Gets the config file from firestore
We only expect there to be one document within this collection

	hence we limit the query to one
*/
func getCollection(ctx context.Context, client *firestore.Client, collection string, w http.ResponseWriter) {
	iter := client.Collection(collection).Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		http.Error(w, "No document found in esp settings", http.StatusNotFound)
	}
	if err != nil {
		http.Error(w, "Error fetching document", http.StatusInternalServerError)
	}

	document := doc.Data()
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(document)
	if err != nil {
		http.Error(w, "Config file format error", http.StatusInternalServerError)
	}

}

func updateCollection(ctx context.Context, client *firestore.Client, collection string, w http.ResponseWriter, r *http.Request) {

	// Query collection
	iter := client.Collection(collection).Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		http.Error(w, "No document found in collection", http.StatusNotFound)
	}
	if err != nil {
		http.Error(w, "Error fetching document in update", http.StatusInternalServerError)
	}

	// Decode payload
	var updateData map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Error decoding body of request", http.StatusBadRequest)
	}

	// Update the document
	_, err = doc.Ref.Set(ctx, updateData, firestore.MergeAll)
	if err != nil {
		http.Error(w, "Error updating", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success"}
	_ = json.NewEncoder(w).Encode(response)

}
