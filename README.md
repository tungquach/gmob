# Gmob

Gmob is simple util for [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver) to quickly build an unordered representation of a BSON document (M) from Map or Struct.

## Install

```bash
go get -u github.com/tungquach/gmob
```

## Usage

```go
package main

import (
    ...
    github.com/tungquach/gmob
    ...
)

// Car example entity
type Car struct {
    ID string `bson:"_id"`
    Name string `bson:"name"`
    CreatedAt time.Time `bson:"createdAt"`
}

func main() {
    // connect to mongodb
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // cars collection
    ctx := context.Background()
    collection := client.Database("test").Collection("cars")

    // find record
    filter := gmob.Build(Car{Name: "BMW"})
    result := &Car{}
    err = collection.FindOne(ctx, filter).Decode(result)
    log.Printf("result %+v, err: %v", result, err)

    // update record
    setValues := gmob.Build(Car{Name: "BMW 2020"})
    err := collection.FindOneAndUpdate(ctx, bson.M{"_id": "car001"},
        bson.M{"$set": setValues}).Err()
    log.Printf("err: %v", err)
}
```