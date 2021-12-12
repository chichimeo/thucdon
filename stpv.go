package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func getAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var temp []Addition
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	cursor, err := collectionAddition.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Addition
		cursor.Decode(&person)
		temp = append(temp, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(temp)
}
