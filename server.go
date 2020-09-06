package main

import (
	"airquality/database/influx"
	"fmt"
	"github.com/btubbs/datetime"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Handles incoming data from sensors, writes it into the database.
// Will fail if any value is missing, corrupted, or not parseable.
func dataHandler(db *influx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Could not parse form stringData", http.StatusBadRequest)
			return
		}

		stringData := map[string]string{}
		stringData["timestampStr"] = r.FormValue("timestamp")         // timestamp of collected reading
		stringData["sensorName"] = r.FormValue("sensorname")          // name of sensor
		stringData["sensorId"] = r.FormValue("sensorid")              // id of sensor
		stringData["particles003dl"] = r.FormValue("particles_003dl") // 00.3 µm per 0.1 liters
		stringData["particles005dl"] = r.FormValue("particles_005dl") // 00.5 µm per 0.1 liters
		stringData["particles010dl"] = r.FormValue("particles_010dl") // 01.0 µm per 0.1 liters
		stringData["particles025dl"] = r.FormValue("particles_025dl") // 02.5 µm per 0.1 liters
		stringData["particles050dl"] = r.FormValue("particles_050dl") // 05.0 µm per 0.1 liters
		stringData["particles100dl"] = r.FormValue("particles_100dl") // 10.0 µm per 0.1 liters
		stringData["stdPm010"] = r.FormValue("stdPm010")              // Standardized readings
		stringData["stdPm025"] = r.FormValue("stdPm025")
		stringData["stdPm100"] = r.FormValue("stdPm100")
		stringData["envPm010"] = r.FormValue("envPm010") // Environmental readings
		stringData["envPm025"] = r.FormValue("envPm025")
		stringData["envPm100"] = r.FormValue("envPm100")

		fmt.Printf("Received form data %s\n", stringData)

		for key, val := range stringData {
			if val == "" {
				http.Error(w, "Missing one or more values, including "+key, http.StatusBadRequest)
				return
			}
		}

		timestamp, dateErr := datetime.Parse(stringData["timestampStr"], time.UTC)
		if dateErr != nil {
			http.Error(w, "Bad date "+stringData["timestampStr"], http.StatusBadRequest)
			return
		}

		// These are the keys of the data that should be kept as strings.
		stringKeys := map[string]bool{"timestampStr": true, "sensorName": true, "sensorId": true}
		// These are the data items that should be floats.
		floatData := map[string]float64{}
		for key, val := range stringData {
			if _, ok := stringKeys[key]; !ok {
				// Cast it into a float
				num, floatCastErr := strconv.ParseFloat(val, 64)
				if floatCastErr == nil {
					floatData[key] = num
				} else {
					http.Error(w, "Bad value "+val+" for key "+key, http.StatusBadRequest)
					return
				}
			}
		}

		// Write data asynchronously
		db.WriteData(timestamp, stringData["sensorName"], stringData["sensorId"], floatData)
	})
}

func main() {
	db := influx.NewDB("http://localhost:8086", "token", "airquality", "airquality")

	http.Handle("/data", dataHandler(db))

	port := os.Getenv("port")
	if port == "" {
		log.Fatal("No port specified")
		return
	}

	fmt.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
