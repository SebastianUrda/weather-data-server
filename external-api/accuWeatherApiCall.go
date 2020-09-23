package external_api

import (
	"awesomeProject/entity"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func InsertAccurateWeatherData(lat float64, lng float64, timestamp time.Time) {
	db, err := sql.Open("mysql", "root:sebi@tcp(127.0.0.1:3306)/licenta?parseTime=true")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	latStr := fmt.Sprintf("%f", lat)
	lngStr := fmt.Sprintf("%f", lng)
	key := "API_KEY"
	locationResp, err := http.Get("http://dataservice.accuweather.com/locations/v1/cities/geoposition/search?apikey=" + key + "&q=" + latStr + "," + lngStr)
	if err != nil {
		log.Fatal(err)
	}
	defer locationResp.Body.Close()
	locationBody, err := ioutil.ReadAll(locationResp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var location entity.AccuWeatherGeoResponse
	json.Unmarshal(locationBody, &location)

	resp, err := http.Get("http://dataservice.accuweather.com/currentconditions/v1/" + location.Key + "?apikey=" + key + "&details=true")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result []*entity.AccuWeatherResponse
	json.Unmarshal(body, &result)
	stmtInsData, err := db.Prepare("INSERT INTO api_observations(latitude, longitude, measurement_unit, timestamp, value, api_name, measuring) VALUES (?,?,?,?,?,?,?)")
	_, err = stmtInsData.Exec(lat, lng, result[0].Temperature.Metric.Unit, timestamp, result[0].Temperature.Metric.Value, "AccuWeather", "Temperature")
	_, err = stmtInsData.Exec(lat, lng, "%", timestamp, result[0].RelativeHumidity, "AccuWeather", "Humidity")
	_, err = stmtInsData.Exec(lat, lng, result[0].Pressure.Metric.Unit, timestamp, result[0].Pressure.Metric.Value, "AccuWeather", "Pressure")
	_, err = stmtInsData.Exec(lat, lng, "", timestamp, result[0].UVIndex, "AccuWeather", "UV Index")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("temperature: %.2f Celsius\n", result[0].Temperature.Metric.Value)
}
