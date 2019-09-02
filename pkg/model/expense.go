package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Expense struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	UserID   string        `bson:"userID" json:"userID"`
	Title    string        `bson:"title" json:"title"`
	Price    float32       `bson:"price" json:"price"`
	Date     string        `bson:"date" json:"date"`
	IsSaving bool          `bson:"isSaving" json:"isSaving"`
	Icon     string        `bson:"icon" json:"icon"`
}
