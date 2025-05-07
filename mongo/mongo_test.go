package mongo

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

const mgoDsn = "mongodb://root:root@127.0.0.1:10000"

var mgoClient *mongo.Client

func init() {
	var err error
	if mgoClient, err = New(mgoDsn); err != nil {
		panic(err)
	}
}

func TestDecimal(t *testing.T) {
	collection := mgoClient.Database("sport").Collection("decimal_test")

	type Card struct {
		Amount decimal.Decimal      `bson:"amount"`
		Money  primitive.Decimal128 `bson:"money"`
	}

	num := "-23245.12345678910"

	amount, _ := decimal.NewFromString(num)
	amount128, _ := primitive.ParseDecimal128(num)
	_, err := collection.InsertOne(context.Background(), Card{Amount: amount, Money: amount128})
	if err != nil {
		panic(err)
	}
	res := collection.FindOne(context.Background(), bson.M{})
	var card Card
	if err = res.Decode(&card); err != nil {
		panic(err)
	}
	fmt.Printf("%+v", card)
}
