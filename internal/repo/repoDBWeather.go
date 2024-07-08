package repo

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/mbatimel/WB_Weather_Service/internal/model"
)

// GetAllCities retrieves a list of all cities sorted by name.
func (db *DataBase) GetAllCities() ([]string, error) {
	var cities []string
	err := db.DB.Model((*model.Cities)(nil)).Column("name").Order("name").Select(&cities)
	if err != nil {
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
			array_agg(w.date ORDER BY w.date) AS available_dates
		FROM 
			cities c
		JOIN 
			weather_forecasts w ON c.id = w.city_id
		WHERE 
			c.name = ?
		GROUP BY 
			c.country, c.name;
	`

	_, err := db.DB.QueryOne(pg.Scan(&result.Country, &result.City, &result.AvgTemp, pg.Array(&result.AvailableDates)), query, cityName)
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
			c.name = ? AND w.date = ?
	`

	_, err := db.DB.QueryOne(pg.Scan(&result.Country, &result.City, &result.Date, &result.Temp, result.WeatherData), query, cityName, date)
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
