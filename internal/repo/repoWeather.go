package repo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"time"

	"github.com/mbatimel/WB_Weather_Service/internal/model"
)

const geocodingAPIURL = "http://api.openweathermap.org/geo/1.0/direct"
const openWeatherAPIURL = "https://api.openweathermap.org/data/2.5/forecast"
const openWeatherAPIKey = "ea97a3b324b49ab2208278142513501d"
func (db *DataBase) InitializeCities() error {

    cityNames := []string{"London", "Paris", "Berlin", "Tokyo"}

    for _, cityName := range cityNames {
        url := fmt.Sprintf("%s?q=%s&limit=1&appid=%s",geocodingAPIURL, cityName, openWeatherAPIKey)
		log.Println(url)
        resp, err := http.Get(url)
        if err != nil {
            return fmt.Errorf("failed to fetch city data: %w", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
        }

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            return fmt.Errorf("failed to read response body: %w", err)
        }

        cityData := make([]map[string]interface{}, 0)
        err = json.Unmarshal(body, &cityData)
        if err != nil {
            return fmt.Errorf("failed to unmarshal city data: %w", err)
        }

        if len(cityData) == 0 {
            log.Printf("no data found for city %s", cityName)
            continue
        }

        city := model.Cities{
            Name:      cityName,
            Country:   cityData[0]["country"].(string),
            Latitude:  cityData[0]["lat"].(float64),
            Longitude: cityData[0]["lon"].(float64),
        }

        _, err = db.DB.Model(&city).Insert()
        if err != nil {
            return fmt.Errorf("failed to insert city %s: %w", city.Name, err)
        }
    }

    return nil
}

func (db *DataBase) UpdateWeatherForecast() error {
    var cities []model.Cities
    err := db.DB.Model(&cities).Select()
    if err != nil {
        return fmt.Errorf("failed to select cities: %w", err)
    }

    for _, city := range cities {
        err := db.updateWeatherForCity(city)
        if err != nil {
            log.Printf("failed to update weather for city %s: %v", city.Name, err)
        }
    }

    return nil
}

func (db *DataBase) updateWeatherForCity(city model.Cities) error {
    url := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s",openWeatherAPIURL, city.Latitude, city.Longitude, openWeatherAPIKey)
	log.Println(url)
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("failed to fetch weather data: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %w", err)
    }

    weatherData := make(map[string]interface{})
    err = json.Unmarshal(body, &weatherData)
    if err != nil {
        return fmt.Errorf("failed to unmarshal weather data: %w", err)
    }

    forecastList, ok := weatherData["list"].([]interface{})
    if !ok {
        return fmt.Errorf("unexpected format for weather data")
    }

    for _, forecast := range forecastList {
        forecastMap := forecast.(map[string]interface{})
        if forecastMap == nil {
            return fmt.Errorf("forecastMap nil")
        }
        mainData := forecastMap["main"].(map[string]interface{})
        if mainData == nil {
            return fmt.Errorf("mainData nil")
        }
		forecastTime := forecastMap["dt_txt"].(string)
        if forecastTime == "" {
            return fmt.Errorf("forecastTime nil")
        }
		
        temp := mainData["temp"].(float64)
        if temp == 0 {
            return fmt.Errorf("temp nil")
        }
        date, err := time.Parse("2006-01-02 15:04:05", forecastMap["dt_txt"].(string))
        if err != nil {
            return fmt.Errorf("failed to parse date: %w", err)
        }

        weatherBytes, err := json.Marshal(forecastMap)
        if err != nil {
            return fmt.Errorf("failed to marshal weather data: %w", err)
        }

        weatherForecast := model.WeatherForecast{
            IdCity:      city.Id,
            Temp:        temp,
            Date:        date,
            WeatherData: weatherBytes,
        }

        _, err = db.DB.Model(&weatherForecast).
            Table("weather_forecasts").
            OnConflict("(city_id, date) DO UPDATE").
            Set("temp = EXCLUDED.temp, weather_data = EXCLUDED.weather_data").
            Insert()
        if err != nil {
            return fmt.Errorf("failed to insert/update weather data: %w", err)
        }
    }

    return nil
}
