package main

import (
	externalapi "externa-api-server/external-api"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/mightWork/insertApiData", handleInsertApiParametersRequest)
	http.HandleFunc("/", replyWelcome)
	http.ListenAndServe(":8080", nil)
}

func handleInsertApiParametersRequest(w http.ResponseWriter, r *http.Request) {

	latitude, ok := r.URL.Query()["latitude"]
	if !ok || len(latitude) < 1 {
		http.Error(w, "Missing latitude", http.StatusInternalServerError)
		log.Println("Missing latitude")
	}

	longitude, ok := r.URL.Query()["longitude"]
	if !ok || len(longitude) < 1 {
		http.Error(w, "Missing longitude", http.StatusInternalServerError)
		log.Println("Missing longitude")
	}
	timestamp, ok := r.URL.Query()["timestamp"]
	if !ok || len(timestamp) < 1 {
		http.Error(w, "Missing timestamp", http.StatusInternalServerError)
		log.Println("Missing timestamp")
	}


	externalapi.InsertDarkSkyData(getFloat64(latitude[0]), getFloat64(longitude[0]),getTimestamp(timestamp[0]))
	externalapi.InsertAccurateWeatherData(getFloat64(latitude[0]), getFloat64(longitude[0]),getTimestamp(timestamp[0]))
}

func replyWelcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func getTimestamp(timestamp string) time.Time {
	//http://localhost:8080/mightWork/insertApiData?latitude=47.8&longitude=22.90&timestamp=2020-06-18T10:52:32Z
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		log.Println(err)
	}
	return t
}

func getFloat64(floatString string) float64{
	floatResponse, err := strconv.ParseFloat(floatString, 64)
	if err != nil {
		log.Println(err)
	}
	return floatResponse
}
