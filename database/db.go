package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DB struct{
	client mongo.Client
}

// TODO
func NewDB() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
}

// TODO
func (db *DB) WriteDataEntry(timestamp time.Time,
	sensorName, sensorId string, datapoints map[string]float64) error {

}
