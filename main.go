package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionMenu *mongo.Collection
var collectionAddition *mongo.Collection

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb+srv://mongo:12345678Abc@datatest.rxmvq.mongodb.net/")
	client, _ := mongo.Connect(ctx, clientOptions)
	collectionMenu = client.Database("bigdata").Collection("menu")
	collectionAddition = client.Database("bigdata").Collection("addition")
	router := mux.NewRouter().StrictSlash(true)
	router.Path("/menu/").Queries("conditions", "{conditions}", "page", "{page}", "limit", "{limit}").HandlerFunc(getMenus).Methods("GET")
	router.HandleFunc("/menu/{foodCode}", getByID).Methods("GET")
	router.HandleFunc("/menu", createMenu).Methods("POST")
	router.HandleFunc("/menu/{foodCode}", updateMenu).Methods("PUT")
	router.HandleFunc("/menu/{foodCode}", deleteMenu).Methods("DELETE")
	router.HandleFunc("/addition", getAll).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", router))
}
