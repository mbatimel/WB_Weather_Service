package repo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/mbatimel/WB_Weather_Service/internal/model"
)

const (
	geocodingAPIURL   = "http://api.openweathermap.org/geo/1.0/direct"
	openWeatherAPIURL = "https://api.openweathermap.org/data/2.5/forecast"
	openWeatherAPIKey = "api"
)

func (db *DataBase) InitializeCities() error {
	cityNames := []string{"London", "Paris", "Berlin", "New York", "Tokyo"}

	for _, cityName := range cityNames {
		url := fmt.Sprintf("%s?q=%s&limit=1&appid=%s", geocodingAPIURL, cityName, openWeatherAPIKey)
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

		// _, err = db.DB.Exec(`INSERT INTO cities (name, country, latitude, longitude) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`, city.Name, city.Country, city.Latitude, city.Longitude)
		// if err != nil {
		// 	return fmt.Errorf("failed to insert city %s: %w", city.Name, err)
		// }
        _, err = db.DB.Model(&city).Insert()
        if err != nil {
            return fmt.Errorf("failed to insert city %s: %w", city.Name, err)
        }
    }



	return nil
}

func (db *DataBase) UpdateWeatherForecast() error {
	var cities []model.Cities
	err := db.DB.Select(&cities)
	if err != nil {
		return fmt.Errorf("failed to select cities: %w", err)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(cities))
	defer close(errCh)

	for _, city := range cities {
		wg.Add(1)
		go func(city model.Cities) {
			defer wg.Done()
			if err := db.updateWeatherForCity(city); err != nil {
				errCh <- fmt.Errorf("failed to update weather for city %s: %v", city.Name, err)
			}
		}(city)
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (db *DataBase) updateWeatherForCity(city model.Cities) error {
	url := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s", openWeatherAPIURL, city.Latitude, city.Longitude, openWeatherAPIKey)
	log.Println("Fetching weather data from:", url)

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
		forecastMap, ok := forecast.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected format for forecast map")
		}

		mainData, ok := forecastMap["main"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected format for main data")
		}

		forecastTime, ok := forecastMap["dt_txt"].(string)
		if !ok {
			return fmt.Errorf("unexpected format for forecast time")
		}
		log.Println(forecastTime)

		temp, ok := mainData["temp"].(float64)
		if !ok {
			return fmt.Errorf("unexpected format for temperature")
		}
		log.Println(temp)

		date, err := time.Parse("2006-01-02 15:04:05", forecastTime)
		if err != nil {
			return fmt.Errorf("failed to parse date: %w", err)
		}
		log.Println(date)

		weatherBytes, err := json.Marshal(forecastMap)
		if err != nil {
			return fmt.Errorf("failed to marshal weather data: %w", err)
		}

	// 	query := `
	// 		INSERT INTO weather_forecasts (city_id, temp, date, weather_data)
	// 		VALUES ($1, $2, $3, $4)
	// 		ON CONFLICT (city_id, date) 
	// 		DO UPDATE SET temp = EXCLUDED.temp, weather_data = EXCLUDED.weather_data;
	// 	`
	// 	_, err = db.DB.Exec(query, city.Id, temp, date, weatherBytes)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to insert/update weather data: %w", err)
	// 	}
	// }

	// return nil
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
