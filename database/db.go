package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

const DB_NAME string = "airquality"
const SENSOR_DATA_COLL string = "sensordata"

type DB struct {
	client               *mongo.Client
	sensorDataCollection *mongo.Collection
}

func NewDB() (*DB, error) {
	// Create the client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println("Error creating mongodb client " + err.Error())
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the database
	if err := client.Connect(ctx); err != nil {
		log.Println("Error connecting to mongodb database " + err.Error())
		return nil, err
	}

	// Check to make sure the server was connected to
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("Error verifying connection to mongodb " + err.Error())
		return nil, err
	}

	newDb := DB{
		client:               client,
		sensorDataCollection: client.Database(DB_NAME).Collection(SENSOR_DATA_COLL),
	}

	return &newDb, nil
}

func (db *DB) WriteDataEntry(timestamp time.Time,
	sensorName, sensorId string, datapoints map[string]float64) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data := bson.M{
		"timestamp":  timestamp,
		"sensorname": sensorName,
		"sensorid":   sensorId,
	}

	for key, val := range datapoints {
		data[key] = val
	}

	_, err := db.sensorDataCollection.InsertOne(ctx, data)
	if err != nil {
		log.Println("Error inserting data entry " + err.Error())
		return err
	}
	return nil
}

func (db *DB) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = db.client.Disconnect(ctx)
}
