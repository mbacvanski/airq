package main

import (
	"airquality/database"
	"fmt"
	"github.com/araddon/dateparse"
	"net/http"
	"strconv"
)

func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not parse form stringData", http.StatusBadRequest)
		return
	}

	stringData := map[string]string{}
	stringData["timestampStr"] = r.FormValue("timestamp")            // timestamp of collected reading
	stringData["sensorName"] = r.FormValue("sensorname")             // name of sensor
	stringData["sensorId"] = r.FormValue("sensorid")                 // id of sensor
	stringData["particles003dlStr"] = r.FormValue("particles_003dl") // 00.3 µm per 0.1 liters
	stringData["particles005dlStr"] = r.FormValue("particles_005dl") // 00.5 µm per 0.1 liters
	stringData["particles010dlStr"] = r.FormValue("particles_010dl") // 01.0 µm per 0.1 liters
	stringData["particles025dlStr"] = r.FormValue("particles_025dl") // 02.5 µm per 0.1 liters
	stringData["particles050dlStr"] = r.FormValue("particles_050dl") // 05.0 µm per 0.1 liters
	stringData["particles100dlStr"] = r.FormValue("particles_100dl") // 10.0 µm per 0.1 liters
	stringData["stdPm010"] = r.FormValue("stdPm010")                 // Standardized readings
	stringData["stdPm025"] = r.FormValue("stdPm025")
	stringData["stdPm100"] = r.FormValue("stdPm100")
	stringData["envPm010"] = r.FormValue("envPm010") // Environmental readings
	stringData["envPm025"] = r.FormValue("envPm025")
	stringData["envPm100"] = r.FormValue("envPm100")

	for key, val := range stringData {
		if val == "" {
			http.Error(w, "Missing one or more values, including "+key, http.StatusBadRequest)
		}
	}

	timestamp, err := dateparse.ParseAny(stringData["timestampStr"])
	if err != nil {
		http.Error(w, "Bad date "+stringData["timestampStr"], http.StatusBadRequest)
	}

	// These are the keys of the data that should be kept as strings.
	stringKeys := map[string]bool{"timestampStr": true, "sensorName": true, "sensorId": true}
	// These are the data items that should be floats.
	floatData := map[string]float64{}
	for key, val := range stringData {
		if _, ok := stringKeys[key]; !ok {
			// Cast it into a float
			num, err := strconv.ParseFloat(val, 64)
			if err == nil {
				floatData[key] = num
			} else {
				http.Error(w, "Bad value "+val+" for key "+key, http.StatusBadRequest)
			}
		}
	}

	// TODO
	newDb := database.NewDB()
	newDb.WriteDataEntry()

	fmt.Printf("Received stringData %s\n")
}

func main() {
	http.Handle("/data", dataHandler)
}
