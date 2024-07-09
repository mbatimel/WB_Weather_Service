package repo

import (
	"context"
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
	openWeatherAPIKey = "ea97a3b324b49ab2208278142513501d"
)

func (db *DataBase) InitializeCities() error {
	cityNames := []string{"London", "Paris", "Berlin", "New York", "Tokyo","Moscow","Saint-Petersburg", "Kazan", "Chelyabinsk", "Novosibirsk","Ekaterinburg","Samara","Omsk","Edinburgh","Cardiff","Belfast","Glasgow", "Manchester","Liverpool","Oslo"}

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
        var exists bool
        err = db.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM cities WHERE name=$1)", cityName).Scan(&exists)
        if err != nil {
            return err
        }
            if !exists {
            query := `INSERT INTO cities (name, country, latitude, longitude) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
            _, err = db.DB.Exec(context.Background(), query, city.Name, city.Country, city.Latitude, city.Longitude)
            if err != nil {
                return fmt.Errorf("InitializeCities: failed to insert city %s: %w", city.Name, err)
            }
        }
	}

	return nil
}
func (db *DataBase) UpdateWeatherForecast() error {
    var cities []model.Cities
    query := `SELECT id, name, country, latitude, longitude FROM cities`
    
    // Execute the query and scan results into cities slice
    rows, err := db.DB.Query(context.Background(), query)
    if err != nil {
        return fmt.Errorf("UpdateWeatherForecast: failed to select cities: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var city model.Cities
        err := rows.Scan(&city.Id, &city.Name, &city.Country, &city.Latitude, &city.Longitude)
        if err != nil {
            return fmt.Errorf("UpdateWeatherForecast: failed to scan city row: %w", err)
        }
        cities = append(cities, city)
    }
    if err := rows.Err(); err != nil {
        return fmt.Errorf("UpdateWeatherForecast: error iterating over city rows: %w", err)
    }

    var wg sync.WaitGroup
    errCh := make(chan error, len(cities))
    defer close(errCh)

	wg.Add(len(cities))
    for _, city := range cities {
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

        temp, ok := mainData["temp"].(float64)
        if !ok {
            return fmt.Errorf("unexpected format for temperature")
        }

        date, err := time.Parse("2006-01-02 15:04:05", forecastTime)
        if err != nil {
            return fmt.Errorf("failed to parse date: %w", err)
        }

        weatherBytes, err := json.Marshal(forecastMap)
        if err != nil {
            return fmt.Errorf("failed to marshal weather data: %w", err)
        }

        query := `
            INSERT INTO weather_forecasts (city_id, temp, date, weather_data)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (city_id, date) 
            DO UPDATE SET temp = EXCLUDED.temp, weather_data = EXCLUDED.weather_data;
        `
        _, err = db.DB.Exec(context.Background(), query, city.Id, temp, date, weatherBytes)
        if err != nil {
            return fmt.Errorf("failed to insert/update weather data: %w", err)
        }
    }

    return nil
}