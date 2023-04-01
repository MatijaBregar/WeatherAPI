package main

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"unicode"
	"weatherAPI/secrets"
	"weatherAPI/utils"
)

func emptyCache() utils.WeatherData {
	weatherData := utils.WeatherData{
		City:     "None",
		Country:  "None",
		Temp_C:   0.0,
		Temp_F:   0.0,
		Wind_KPH: 0.0,
		Wind_MPH: 0.0,
	}
	return weatherData
}

func isNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func getWeatherByCity(city string) (*utils.WeatherData, error) {
	url := "https://weatherapi-com.p.rapidapi.com/current.json?q=" + city

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", secrets.GetAPIKey())
	req.Header.Add("X-RapidAPI-Host", "weatherapi-com.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, errors.New("Service is not available")
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var data map[string]interface{}
	err2 := json.Unmarshal([]byte(string(body)), &data)

	if err2 != nil {
		return nil, errors.New("Error while parsing data")
	}

	var weatherData utils.WeatherData

	if location, ok := data["location"].(map[string]interface{}); ok {
		if current, ok := data["current"].(map[string]interface{}); ok {
			weatherData := utils.WeatherData{
				City:     location["name"].(string),
				Country:  location["country"].(string),
				Temp_C:   current["temp_c"].(float64),
				Temp_F:   current["temp_f"].(float64),
				Wind_KPH: current["wind_kph"].(float64),
				Wind_MPH: current["wind_mph"].(float64),
				Is_day:   current["is_day"].(float64),
				Wind_dir: current["wind_dir"].(string),
				Humidity: current["humidity"].(float64),
			}
			keyValue := utils.KeyValue{Key: city, Value: weatherData}
			addToCache(keyValue)
			return &weatherData, nil
		} else {
			return nil, errors.New("Current informations not found")
		}
	} else {
		if isNumeric(city) {
			return nil, errors.New("Zip code not found")
		}
		return nil, errors.New("Location not found")
	}

	return &weatherData, nil
}

func getWeather(context *gin.Context) {
	city := context.Param("city")

	rezult, err := fetchFromCache(city)

	if err != nil {
		weatherData, err := getWeatherByCity(city)
		if err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		context.IndentedJSON(http.StatusOK, weatherData)
	} else {
		context.IndentedJSON(http.StatusOK, rezult)
	}
}

func addToCache(object utils.KeyValue) error {
	err := utils.CacheStore.Add(object)

	if err != nil {
		klog.Errorf("failed to add key value to cache error", err)
		return err
	}
	return nil
}

func fetchFromCache(key string) (utils.WeatherData, error) {
	obj, exists, err := utils.CacheStore.GetByKey(key)

	if err != nil {
		klog.Errorf("failed to add key value to cache error", err)
		weatherData := emptyCache()
		return weatherData, err
	}
	if !exists {
		klog.Errorf("object does not exist in the cache")
		weatherData := emptyCache()
		return weatherData, errors.New("an error occurred")
	}

	klog.Errorf("Object found in cache")
	return obj.(utils.KeyValue).Value, nil
}

func deleteFromCache(object utils.KeyValue) error {
	return utils.CacheStore.Delete(object)
}

func main() {
	r := gin.Default()
	r.GET("/:city", getWeather)
	r.Run("localhost:8080")
}
