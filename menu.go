package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Datatable(items interface{}, filter, p, l string, searchable []string) (datatable interface{}, total int, err error) {
	var page, limit int

	if p == "" {
		p = "1"
	}

	if page, err = strconv.Atoi(p); err != nil {
		return
	}

	if l == "" {
		l = "10"
	}

	if limit, err = strconv.Atoi(l); err != nil {
		return
	}

	bsonfilter := bson.M{}
	or := []bson.M{}
	for _, field := range searchable {
		or = append(or, bson.M{field: primitive.Regex{Pattern: filter, Options: "gi"}})
	}

	if len(or) > 0 {
		bsonfilter = bson.M{"$or": or}
	}

	// find record of one page
	cursor, err := collectionMenu.Find(context.Background(),
		bsonfilter,
		options.Find().SetLimit(int64(limit)).SetSkip(int64((page-1)*limit)))
	if err != nil {
		return
	}

	err = cursor.All(context.Background(), &items)
	if err != nil {
		return
	}

	datatable = items

	// total record for pagination
	total64, err := collectionMenu.CountDocuments(context.Background(),
		bsonfilter,
		options.Count())
	if err != nil {
		return
	}

	total = int(total64)
	return
}

func createMenu(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var menu MenuCreate
	json.NewDecoder(request.Body).Decode(&menu)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var resultFindOne Menu
	_ = collectionMenu.FindOne(ctx, bson.M{"foodCode": menu.FoodCode}).Decode(&resultFindOne)
	if resultFindOne.FoodCode == menu.FoodCode {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": " foodCode is duplicate" }`))
		return
	}

	now := time.Now()
	newData := Menu{
		FoodCode:                  menu.FoodCode,
		FoodName:                  menu.FoodName,
		InventoryItemCategoryName: menu.InventoryItemCategoryName,
		Unit:                      menu.Unit,
		SalePrice:                 menu.SalePrice,
		RealPrice:                 menu.RealPrice,
		Description:               menu.Description,
		Kitchen:                   menu.Kitchen,
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}

	defer cancel()
	for _, value := range menu.FoodAdditions {
		if value.ID == "" {
			value.ID = uuid.NewString()
			_, err := collectionAddition.InsertOne(ctx, value)

			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
				return
			}
		}
		newData.FoodAddition = append(newData.FoodAddition, value.ID)
	}
	_, err := collectionMenu.InsertOne(ctx, newData)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(newData)
}

func getMenus(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var temp []Menu
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	filter := request.FormValue("filter")
	page := request.FormValue("page")
	limit := request.FormValue("limit")
	defer cancel()
	searchable := []string{"foodCode", "foodName", "unit"}
	datatable, total, err := Datatable(temp, filter, page, limit, searchable)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	table, ok := datatable.([]Menu)
	if !ok {
		err = fmt.Errorf("can't cast")
		return
	}

	temp = make([]Menu, len(table))
	for i, item := range table {
		temp[i] = item
	}

	json.NewEncoder(response).Encode(&DataTable{Data: temp, Total: total})
}

func updateMenu(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var menu, resultFindOne MenuCreate
	params := mux.Vars(request)
	foodCode := params["foodCode"]

	json.NewDecoder(request.Body).Decode(&menu)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	err := collectionMenu.FindOne(ctx, bson.M{"foodCode": foodCode}).Decode(&resultFindOne)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	defer cancel()
	now := time.Now()
	newData := Menu{
		FoodCode:                  menu.FoodCode,
		FoodName:                  menu.FoodName,
		InventoryItemCategoryName: menu.InventoryItemCategoryName,
		Unit:                      menu.Unit,
		SalePrice:                 menu.SalePrice,
		RealPrice:                 menu.RealPrice,
		Description:               menu.Description,
		Kitchen:                   menu.Kitchen,
		CreatedAt:                 menu.CreatedAt,
		UpdatedAt:                 now,
	}
	for _, value := range menu.FoodAdditions {
		if value.ID == "" {
			value.ID = uuid.NewString()
			_, err := collectionAddition.InsertOne(ctx, value)

			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
				return
			}
		}
		newData.FoodAddition = append(newData.FoodAddition, value.ID)
	}
	menu.UpdatedAt = now
	result, err := collectionMenu.UpdateOne(ctx, bson.M{"foodCode": foodCode}, bson.M{"$set": newData})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(result)
}

func deleteMenu(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	foodCode, _ := mux.Vars(request)["foodCode"]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := collectionMenu.DeleteOne(ctx, bson.M{"foodCode": foodCode})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}
func getByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var menu Menu
	params := mux.Vars(request)
	foodCode := params["foodCode"]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collectionMenu.FindOne(ctx, bson.M{"foodCode": foodCode}).Decode(&menu)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	cursor, err := collectionAddition.Find(ctx, bson.M{"id": bson.M{"$in": menu.FoodAddition}})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	additions := make([]Addition, len(menu.FoodAddition))
	err = cursor.All(context.Background(), &additions)

	newData := MenuCreate{
		FoodCode:                  menu.FoodCode,
		FoodName:                  menu.FoodName,
		InventoryItemCategoryName: menu.InventoryItemCategoryName,
		Unit:                      menu.Unit,
		SalePrice:                 menu.SalePrice,
		RealPrice:                 menu.RealPrice,
		Description:               menu.Description,
		Kitchen:                   menu.Kitchen,
		CreatedAt:                 menu.CreatedAt,
		UpdatedAt:                 menu.UpdatedAt,
		FoodAdditions:             additions,
	}
	json.NewEncoder(response).Encode(newData)
}
