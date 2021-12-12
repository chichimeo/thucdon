package main

import (
	"time"
)

// Menu represents a menu document in MongoDB
type Menu struct {
	FoodCode                  string    `json:"foodCode" bson:"foodCode"`
	FoodName                  string    `json:"foodName" bson:"foodName"`
	InventoryItemCategoryName string    `json:"inventoryItemCategoryName,omitempty" typeName:"inventoryItemCategoryName,omitempty"`
	Unit                      string    `json:"unit" bson:"unit"`
	SalePrice                 float64   `json:"salePrice" bson:"salePrice"`
	RealPrice                 float64   `json:"realPrice,omitempty" bson:"realPrice,omitempty"`
	Description               string    `json:"description,omitempty" bson:"description,omitempty"`
	Kitchen                   string    `json:"kitchen,omitempty" bson:"kitchen,omitempty"`
	CreatedAt                 time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt                 time.Time `json:"updatedAt" bson:"updatedAt"`
	FoodAddition              []string  `json:"foodAddition" bson:"foodAddition"`
}
type MenuCreate struct {
	FoodCode                  string     `json:"foodCode" bson:"foodCode"`
	FoodName                  string     `json:"foodName" bson:"foodName"`
	InventoryItemCategoryName string     `json:"inventoryItemCategoryName,omitempty" typeName:"inventoryItemCategoryName,omitempty"`
	Unit                      string     `json:"unit" bson:"unit"`
	SalePrice                 float64    `json:"salePrice" bson:"salePrice"`
	RealPrice                 float64    `json:"realPrice,omitempty" bson:"realPrice,omitempty"`
	Description               string     `json:"description,omitempty" bson:"description,omitempty"`
	Kitchen                   string     `json:"kitchen,omitempty" bson:"kitchen,omitempty"`
	CreatedAt                 time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt                 time.Time  `json:"updatedAt" bson:"updatedAt"`
	FoodAdditions             []Addition `json:"foodAdditions" bson:"foodAdditions"`
}
type Addition struct {
	ID    string `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

type DataTable struct {
	Data  []Menu `json:"data" bson:"data"`
	Total int    `json:"total" bson:"total"`
}

type Condition struct {
	Field    string `json:"field" bson:"field"`
	Value    string `json:"value" bson:"value"`
	Type     string `json:"type" bson:"type"`
	Operator string `json:"operator" bson:"operator"`
}
