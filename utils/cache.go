package utils

import (
	"k8s.io/client-go/tools/cache"
	"time"
)

var CacheStore cache.Store

type WeatherData struct {
	City     string  `json:"city"`
	Country  string  `json:"country"`
	Temp_C   float64 `json:"temp_c"`
	Temp_F   float64 `json:"temp_f"`
	Wind_KPH float64 `json:"wind_kph"`
	Wind_MPH float64 `json:"wind_mph"`
	Is_day   float64 `json:"is_day"`
	Wind_dir string  `json:"wind_dir"`
	Humidity float64 `json:"humidity"`
}

type KeyValue struct {
	Key   string
	Value WeatherData
}

func cacheKeyFunc(obj interface{}) (string, error) {
	return obj.(KeyValue).Key, nil
}

func init() {
	cacheTTL := 60 * time.Minute
	CacheStore = cache.NewTTLStore(cacheKeyFunc, cacheTTL)
}
