package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect(outcli **mongo.Client, url string) error {
	opt := options.Client().ApplyURI(url)
	maxConnIdleTime := time.Hour
	opt.MaxConnIdleTime = &maxConnIdleTime
	maxPollSize := uint64(300)
	opt.MaxPoolSize = &maxPollSize
	opt.SetSocketTimeout(time.Minute * 10)
	client, err := mongo.NewClient(opt)
	if err != nil {
		return fmt.Errorf("mongo.NewClient:%v", err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	*outcli = client
	return nil
}

func main() {
	var cli *mongo.Client
	if err := connect(&cli, os.Args[1]); err != nil {
		panic(err)
	}
	fmt.Println(cli)
	result, err := cli.Database("test").Collection("test_coll").InsertOne(nil, bson.M{"a": 1})
	fmt.Println(result, err)
}
