package main

import (
	"os"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"fmt"
	"time"
)

type WeatherData struct {
	City string `json:"city"`
	Info string `json:"info"`
}

var apiKey = os.Getenv("bc94089fef9d46968ff8d724889f9fff")

func main() {
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := r.URL.Path[len("/weather/"):]

		url := "https://weatherapi-com.p.rapidapi.com/current.json?q=" + city

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("X-RapidAPI-Key", "356fcc1b96msh6b0261faa77fcd2p101d17jsnbbb68ac024b7")
		req.Header.Add("X-RapidAPI-Host", "weatherapi-com.p.rapidapi.com")

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		weatherData := WeatherData{
			City:    city,
			Info:    string(body),
		}
		jsonData, _ := json.Marshal(weatherData)
		ioutil.WriteFile("weather.json", jsonData, 0600)

		data, err := ioutil.ReadFile("weather.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, string(data))
	})

	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 3 * time.Second,
	    }

	    err := server.ListenAndServe()
	    if err != nil {
		panic(err)
	    }
}
