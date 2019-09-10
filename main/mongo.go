package main

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Base struct {
	A string
	B string
	c string
}

type Object struct {
	Base
	D string
	e string
}

type Object2 struct {
	A string
	B string
	c string
	D string
	e string
}

func dump(i interface{}) {
	jsonBytes, _ := json.MarshalIndent(i, "", "    ")
	fmt.Println(string(jsonBytes))
}

func main() {

	//recover()

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongodb:27017"))

	if err != nil {
		panic(err)
	}

	err = client.Connect(nil)
	if err != nil {
		panic(err)
	}

	err = client.Ping(nil, nil)
	if err != nil {
		panic(err)
	}
	//insertTest(client)

	{

		cursor, err := client.Database("test").Collection("test1").
			Find(nil, bson.M{}, options.Find().SetLimit(10).SetProjection(bson.M{"a": 1, "_id": 0}))
		if err != nil {
			panic(err)
		}
		for cursor.Next(nil) {
			var o map[string]interface{}
			cursor.Decode(&o)
			fmt.Println(o)
			dump(o)
		}

		d := bson.D{{"x", "x"}}

	}
	return
	{
		r, err := client.Database("test").Collection("test1").DeleteMany(nil, bson.M{})
		if err != nil {
			panic(err)
		}
		fmt.Println(r.DeletedCount)
	}
	{
		o := new(Object)
		o.A = "A"
		o.B = "B"
		o.c = "c"
		o.D = "D"
		o.e = "e"
		r, err := client.Database("test").Collection("test1").InsertOne(nil, o)
		if err != nil {
			panic(err)
		}
		fmt.Println(r.InsertedID)
	}
	{
		var o Object
		err := client.Database("test").Collection("test1").FindOne(nil, bson.M{}).Decode(&o)
		if err != nil {
			panic(err)
		}
		fmt.Println(o)
		dump(o)
	}
	{
		var o Object2
		err := client.Database("test").Collection("test1").FindOne(nil, bson.M{}).Decode(&o)
		if err != nil {
			panic(err)
		}
		fmt.Println(o)
		dump(o)
	}
}

func insertTest(client *mongo.Client) {
	collection := client.Database("test").Collection("test1")
	waitGroup := &sync.WaitGroup{}

	insert := func() {
		_, err := collection.InsertOne(nil, bson.M{fmt.Sprintf("%d", rand.Int()): rand.Int()})

		if err != nil {
			panic(err)
		}
		waitGroup.Done()
		//fmt.Println(result.InsertedID)
	}

	startTime := time.Now()
	n := 10000
	fmt.Println("开始插入测试")
	for i := 0; i < n; i++ {
		waitGroup.Add(1)
		go insert()
	}
	waitGroup.Wait()
	fmt.Printf("插入 %d条， 共花费%s\n", n, time.Now().Sub(startTime))

}
