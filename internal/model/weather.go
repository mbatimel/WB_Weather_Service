package model

type WeatherResponse struct {
	List []WeatherUnit `json:"list"`
}

type WeatherUnit struct {
	Dt         int64           `json:"dt"`
	Main       WeatherMainData `json:"main"`
	Weather    []WeatherWeather  `json:"weather"`
	Clouds     Clouds          `json:"clouds"`
	Wind       Wind            `json:"wind"`
	Visibility float64         `json:"visibility"`
	Pop        float64         `json:"pop"`
	Sys			Sys `json:"sys"`
	DtTxt string `json:"dt_txt"`
}

type WeatherMainData struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  float64 `json:"pressure"`
	SeaLevel  float64 `json:"sea_level"`
	GrndLevel float64 `json:"grnd_level"`
	Humidity  float64 `json:"humidity"`
	TempKF    float64 `json:"temp_kf"`
}

type WeatherWeather struct {
	ID          int64  `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Clouds struct {
	All float64 `json:"all"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   float64 `json:"deg"`
	Gust  float64 `json:"gust"`
}
type Sys struct{
	Pod string `json:"pod"`
}
