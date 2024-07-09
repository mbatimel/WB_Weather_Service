package repo

import (
	"testing"
	"time"
)

func TestInitializeCities(t *testing.T){
	db, err := SetConfigs("../../config/config.yaml")
    want := error(nil)
	if want != err{
		t.Errorf("got %q, wanted %q", err, want)
	}
    defer db.Close()

    err = db.ConnectToDataBase()
    if err != want {
        t.Errorf("got %q, wanted %q", err, want)
    }
	db.InitializeCities()
	wantRes:=[]string{"London", "Paris", "Berlin", "New York", "Tokyo"}
	result, errDB := db.GetAllCities()
	if errDB != want {
		t.Errorf("got %q, wanted %q", errDB, want)
	}
	if len(wantRes) !=len(result){
		t.Errorf("got %q, wanted %q", result, wantRes)
	}
	errUP := db.UpdateWeatherForecast()
	if errUP != want {
		t.Errorf("got %q, wanted %q", errDB, want)
	}
}

func TestUpdateWeatherForecast(t *testing.T){
	db, err := SetConfigs("../../config/config.yaml")
    want := error(nil)
	if want != err{
		t.Errorf("got %q, wanted %q", err, want)
	}
    defer db.Close()

    err = db.ConnectToDataBase()
    if err != want {
        t.Errorf("got %q, wanted %q", err, want)
    }
}
func TestConnect(t *testing.T){
	db, err := SetConfigs("../../config/config.yaml")
    want := error(nil)
	if want != err{
		t.Errorf("got %q, wanted %q", err, want)
	}
    defer db.Close()

    err = db.ConnectToDataBase()
    if err != want {
        t.Errorf("got %q, wanted %q", err, want)
    }
}
func TestGetAllCities(t *testing.T) {
	db, err := SetConfigs("../../config/config.yaml")
    want := error(nil)
	if want != err{
		t.Errorf("got %q, wanted %q", err, want)
	}
    defer db.Close()

    err = db.ConnectToDataBase()
    if err != want {
        t.Errorf("got %q, wanted %q", err, want)
    }
	wantRes:=[]string{"London", "Paris", "Berlin", "New York", "Tokyo"}
	result, errDB := db.GetAllCities()
	if errDB != want {
		t.Errorf("got %q, wanted %q", errDB, want)
	}
	if len(wantRes) != len(result){
		t.Errorf("got %q, wanted %q", result, wantRes)
	}
}

func TestGetShortInfoCity(t *testing.T) {
	db, err := SetConfigs("../../config/config.yaml")
	if err != nil {
		t.Fatalf("Failed to set configs: %v", err)
	}
	defer db.Close()

	err = db.ConnectToDataBase()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	wantRes := map[string]interface{}{
		"country":         "GB",
		"city":            "London",
		"avg_temp":        289.44666666666666, // Example average temperature
		"available_dates": []string{"2024-07-09", "2024-07-10", "2024-07-11", "2024-07-12", "2024-07-13", "2024-07-14"}, // Example dates
	}

	result, err := db.GetShortInfoCity("London")
	if err != nil {
		t.Fatalf("Failed to get short info for city: %v", err)
	}

	if result["country"] != wantRes["country"] || result["city"] != wantRes["city"] || result["avg_temp"] != wantRes["avg_temp"] {
		t.Errorf("got %v, wanted %v", result, wantRes)
	}

	if len(result["available_dates"].([]string)) != len(wantRes["available_dates"].([]string)) {
		t.Errorf("got %v, wanted %v", result["available_dates"], wantRes["available_dates"])
	}
}

func TestGetFullInfoCity(t *testing.T) {
	db, err := SetConfigs("../../config/config.yaml")
	if err != nil {
		t.Fatalf("Failed to set configs: %v", err)
	}
	defer db.Close()

	err = db.ConnectToDataBase()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	wantRes := map[string]interface{}{
		"country": "DE",
		"city":    "Berlin",
		"date":    time.Date(2024, 7, 9, 0, 0, 0, 0, time.UTC),
		"temp":    298.66,
	}

	result, err := db.GetFullInfoCity("Berlin", time.Date(2024, 7, 9, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Failed to get full info for city: %v", err)
	}

	if result["country"] != wantRes["country"] || result["city"] != wantRes["city"] || result["temp"] != wantRes["temp"] {
		t.Errorf("got %v, wanted %v", result, wantRes)
	}

	if !result["date"].(time.Time).Equal(wantRes["date"].(time.Time)) {
		t.Errorf("got %v, wanted %v", result["date"], wantRes["date"])
	}

}
