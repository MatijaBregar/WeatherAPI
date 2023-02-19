package main

import (
	"net/http"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"encoding/json"
)

type WeatherData struct {
	City string `json:"city"`
	Country string `json:"country"`
	Temp_C float64 `json:"temp_c"`
	Temp_F float64 `json:"temp_f"`
	Wind_KPH float64 `json:"wind_kph"`
	Wind_MPH float64 `json:"wind_mph"`
}

func getWeatherByCity(city string) (*WeatherData, error) {
	url := "https://weatherapi-com.p.rapidapi.com/current.json?q=" + city

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "356fcc1b96msh6b0261faa77fcd2p101d17jsnbbb68ac024b7")
	req.Header.Add("X-RapidAPI-Host", "weatherapi-com.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, errors.New("Error getting weather data")
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var data map[string]interface{}
    err2 := json.Unmarshal([]byte(string(body)), &data)

    if err2 != nil {
        return nil, errors.New("Error getting weather data")
    }

	var weatherData WeatherData

	if location, ok := data["location"].(map[string]interface{}); ok {
		if current, ok := data["current"].(map[string]interface{}); ok {
			weatherData := WeatherData{
				City: location["name"].(string),
				Country: location["country"].(string),
				Temp_C: current["temp_c"].(float64),
				Temp_F: current["temp_f"].(float64),
				Wind_KPH: current["wind_kph"].(float64),
				Wind_MPH: current["wind_mph"].(float64),
			}
			return &weatherData, nil
		} else {
			return nil, errors.New("Error getting weather data")
		}
	} else {
		return nil, errors.New("Error getting weather data")
	}

	return &weatherData, nil
}

func getWeather(context *gin.Context) {
	city := context.Param("city")
	weatherData, err := getWeatherByCity(city)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"error": "Was not able to get weather data"})
		return
	}

	context.IndentedJSON(http.StatusOK, weatherData)
}

func main() {
	r := gin.Default()
	r.GET("/:city", getWeather)
	r.Run("localhost:8080")
}