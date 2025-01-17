package model

import (
	"encoding/json"
	"time"
)
type Cities struct{
	Id			int		`json:"city_id" pg:"id,pk"`
	Name		string	`json:"name"`
	Country		string	`json:"country"`
	Latitude	float64	`json:"latitude"`
	Longitude	float64	`json:"longitude"`

}
type WeatherForecast struct{
	Id					int						`json:"weather_id"`
	IdCity				int						`json:"city_id"`
	Temp				float64					`json:"temp"`
	Date				time.Time				`json:"date"`
	WeatherData			json.RawMessage			`json:"weather_data"`
}

type Person struct{
	Id			int			`json:"person_id"`
	Name 		string		`json:"person_name"`
	Password	string		`json:"password"`
	IdCity		[]int		`json:"city_id"`
}
