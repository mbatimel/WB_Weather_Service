package repo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)
func (db *DataBase) AddUser(personName, personPswd string) error {
	// Проверка наличия пользователя с таким именем
	var exists bool
	err := db.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM persons WHERE name=$1)", personName).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(personPswd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Вставка нового пользователя в базу данных
	_, err = db.DB.Exec(context.Background(), "INSERT INTO persons (name, password) VALUES ($1, $2)", personName, string(hashedPassword))
	return err
}
func (db *DataBase) AddCityInFavorit(cityName, personName, personPswd string) error {
	// Проверка существования пользователя и пароля
	var personID int
	var hashedPassword string
	err := db.DB.QueryRow(context.Background(), "SELECT id, password FROM persons WHERE name=$1", personName).Scan(&personID, &hashedPassword)
	if err != nil {
		return err
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(personPswd))
	if err != nil {
		return err
	}

	// Получение id города
	var cityID int
	err = db.DB.QueryRow(context.Background(), "SELECT id FROM cities WHERE name=$1", cityName).Scan(&cityID)
	if err != nil {
		return err
	}

	// Добавление города в избранное
	_, err = db.DB.Exec(context.Background(), "INSERT INTO favorite_cities (person_id, city_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", personID, cityID)
	if err != nil {
		return err
	}

	return nil
}

func (db *DataBase) GetCityInFavorit(personName, personPswd string) (map[string]interface{}, error) {
	// Check if the user exists and verify the password
	var personID int
	var hashedPassword string
	err := db.DB.QueryRow(context.Background(), "SELECT id, password FROM persons WHERE name=$1", personName).Scan(&personID, &hashedPassword)
	if err != nil {
		return nil, err
	}

	// Check the password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(personPswd))
	if err != nil {
		return nil, err
	}

	// Get the list of favorite cities along with weather forecast
	query := `
		SELECT 
			c.name, c.country, c.latitude, c.longitude,
			wf.temp, wf.date, wf.weather_data
		FROM 
			favorite_cities fc
		JOIN 
			cities c ON fc.city_id = c.id
		LEFT JOIN 
			weather_forecasts wf ON c.id = wf.city_id
		WHERE 
			fc.person_id = $1
	`
	rows, err := db.DB.Query(context.Background(), query, personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cities := make([]map[string]interface{}, 0)
	for rows.Next() {
		var name, country string
		var latitude, longitude, temp float64
		var date time.Time
		var weatherData []byte

		if err := rows.Scan(&name, &country, &latitude, &longitude, &temp, &date, &weatherData); err != nil {
			return nil, err
		}

		weather := make(map[string]interface{})
		if err := json.Unmarshal(weatherData, &weather); err != nil {
			return nil, err
		}

		city := map[string]interface{}{
			"name":      name,
			"country":   country,
			"latitude":  latitude,
			"longitude": longitude,
			"weather": map[string]interface{}{
				"temp":         temp,
				"date":         date.Format("2006-01-02"),
				"weather_data": weather,
			},
		}
		cities = append(cities, city)
	}

	return map[string]interface{}{
		"person": personName,
		"cities": cities,
	}, nil
}