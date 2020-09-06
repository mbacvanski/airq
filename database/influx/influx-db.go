package influx

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"time"
)

type DB struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
}

func NewDB(url, token, org, bucket string) *DB {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient(url, token)
	// Get non-blocking write client
	writeAPI := client.WriteAPI(org, bucket)
	// Get errors channel
	errorsCh := writeAPI.Errors()
	// Create go proc for reading and logging errors
	go func() {
		for err := range errorsCh {
			fmt.Printf("write error: %s\n", err.Error())
		}
	}()

	return &DB{
		client:   client,
		writeAPI: writeAPI,
	}
}

func (db DB) WriteData(timestamp time.Time, sensorName, sensorId string, datapoints map[string]float64) {
	tags := map[string]string{
		"sensorName": sensorName,
		"sensorId":   sensorId,
	}

	fields := make(map[string]interface{}, len(datapoints))
	for k, v := range datapoints {
		fields[k] = v
	}

	point := influxdb2.NewPoint(
		"airquality", tags, fields, timestamp)

	// write asynchronously
	db.writeAPI.WritePoint(point)
}

func (db DB) Close() {
	// Force all unwritten data to be sent
	db.writeAPI.Flush()
	// Ensures background processes finishes
	db.client.Close()
}
