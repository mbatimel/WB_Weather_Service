package repo

import (
	"context"
	"time"
)

// GetAllCities retrieves a list of all cities sorted by name.
func (db *DataBase) GetAllCities() ([]string, error) {
	query := `SELECT name FROM cities ORDER BY name`
	rows, err := db.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []string
	for rows.Next() {
		var city string
		if err := rows.Scan(&city); err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cities, nil
}


// GetShortInfoCity retrieves short info for a given city, including average temperature and available forecast dates.
func (db *DataBase) GetShortInfoCity(cityName string) (map[string]interface{}, error) {
	var result struct {
		Country        string
		City           string
		AvgTemp        float64
		AvailableDates []string
	}

	query := `
		SELECT 
			c.country, 
			c.name AS city, 
			AVG(w.temp) AS avg_temp, 
			array_agg(w.date::text ORDER BY w.date) AS available_dates
		FROM 
			cities c
		JOIN 
			weather_forecasts w ON c.id = w.city_id
		WHERE 
			c.name = $1
		GROUP BY 
			c.country, c.name;
	`

	err := db.DB.QueryRow(context.Background(), query, cityName).Scan(&result.Country, &result.City, &result.AvgTemp, &result.AvailableDates)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"country":         result.Country,
		"city":            result.City,
		"avg_temp":        result.AvgTemp,
		"available_dates": result.AvailableDates,
	}, nil
}



// GetFullInfoCity retrieves full weather info for a given city and date.
func (db *DataBase) GetFullInfoCity(cityName string, date time.Time) (map[string]interface{}, error) {
	var result struct {
		Country     string                 `json:"country"`
		City        string                 `json:"city"`
		Date        time.Time              `json:"date"`
		Temp        float64                `json:"temp"`
		WeatherData map[string]interface{} `json:"weather_data"`
	}

	query := `
		SELECT 
			c.country, 
			c.name AS city, 
			w.date, 
			w.temp, 
			w.weather_data 
		FROM 
			cities c
		JOIN 
			weather_forecasts w ON c.id = w.city_id
		WHERE 
			c.name = $1 AND w.date = $2
	`

	err := db.DB.QueryRow(context.Background(), query, cityName, date).Scan(&result.Country, &result.City, &result.Date, &result.Temp, &result.WeatherData)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"country":      result.Country,
		"city":         result.City,
		"date":         result.Date,
		"temp":         result.Temp,
		"weather_data": result.WeatherData,
	}, nil
}

