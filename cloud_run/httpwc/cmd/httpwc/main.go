package main

import (
	"context"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/daveszczesny/web-connected-chess-clock/internal/timecontrol"
)

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

	http.HandleFunc("/timecontrol", timecontrol.HandleCollection(ctx, client, "time-control"))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Could not listen")
	}
}
