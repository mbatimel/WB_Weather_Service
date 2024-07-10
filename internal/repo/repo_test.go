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
	wantRes:=[]string{"London", "Paris", "Berlin", "New York", "Tokyo","Moscow","Saint-Petersburg", "Kazan", "Chelyabinsk", "Novosibirsk","Ekaterinburg","Samara","Omsk","Edinburgh","Cardiff","Belfast","Glasgow", "Manchester","Liverpool","Oslo"}
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
	wantRes:=[]string{"London", "Paris", "Berlin", "New York", "Tokyo","Moscow","Saint-Petersburg", "Kazan", "Chelyabinsk", "Novosibirsk","Ekaterinburg","Samara","Omsk","Edinburgh","Cardiff","Belfast","Glasgow", "Manchester","Liverpool","Oslo"}
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
		

	}

	result, err := db.GetShortInfoCity("London")
	if err != nil {
		t.Fatalf("Failed to get short info for city: %v", err)
	}

	if result["country"] != wantRes["country"] || result["city"] != wantRes["city"] {
		t.Errorf("got %v, wanted %v", result, wantRes)
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
	currentDate := time.Now().Truncate(24 * time.Hour)
	wantRes := map[string]interface{}{
		"country": "DE",
		"city":    "Berlin",
		"date":    currentDate,
	}

	result, err := db.GetFullInfoCity("Berlin", currentDate)
	if err != nil {
		t.Fatalf("Failed to get full info for city: %v", err)
	}

	if result["country"] != wantRes["country"] || result["city"] != wantRes["city"]{
		t.Errorf("got %v, wanted %v", result, wantRes)
	}

	if !result["date"].(time.Time).Equal(wantRes["date"].(time.Time)) {
		t.Errorf("got %v, wanted %v", result["date"], wantRes["date"])
	}

}
