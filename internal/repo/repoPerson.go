package repo

import (
	"context"
	"errors"

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
	// Проверка существования пользователя и пароля
	var personID int
	var hashedPassword string
	err := db.DB.QueryRow(context.Background(), "SELECT id, password FROM persons WHERE name=$1", personName).Scan(&personID, &hashedPassword)
	if err != nil {
		return nil, err
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(personPswd))
	if err != nil {
		return nil, err
	}

	// Получение списка избранных городов
	rows, err := db.DB.Query(context.Background(), "SELECT c.name, c.country, c.latitude, c.longitude FROM favorite_cities fc JOIN cities c ON fc.city_id = c.id WHERE fc.person_id = $1", personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cities := make([]map[string]interface{}, 0)
	for rows.Next() {
		var name, country string
		var latitude, longitude float64
		if err := rows.Scan(&name, &country, &latitude, &longitude); err != nil {
			return nil, err
		}
		city := map[string]interface{}{
			"name":      name,
			"country":   country,
			"latitude":  latitude,
			"longitude": longitude,
		}
		cities = append(cities, city)
	}

	return map[string]interface{}{
		"person": personName,
		"cities": cities,
	}, nil
}