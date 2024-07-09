package repo

import (
	"testing"
)

func TestInitializeCities(t *testing.T){
	db, err := SetConfigs("config/config.yaml")
    want := error(nil)
	if want != err{
		t.Errorf("got %q, wanted %q", err, want)
	}
    defer db.Close()

    err = db.ConnectToDataBase()
    if err != want {
        t.Errorf("got %q, wanted %q", err, want)
    }
	//посмотреть что в БД закинуто

}

func TestUpdateWeatherForecast(t *testing.T){
	db, err := SetConfigs("config/config.yaml")
    want := error(nil)
	if want != err{
		t.Errorf("got %q, wanted %q", err, want)
	}
    defer db.Close()

    err = db.ConnectToDataBase()
    if err != want {
        t.Errorf("got %q, wanted %q", err, want)
    }
	//посмотреть что в БД закинуто
}
func TestConnect(t *testing.T){
	db, err := SetConfigs("config/config.yaml")
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
	db, err := SetConfigs("config/config.yaml")
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
	if len(wantRes) !=len(result){
		t.Errorf("got %q, wanted %q", result, wantRes)
	}
}

func TestGetShortInfoCity(t *testing.T) {
	db, err := SetConfigs("config/config.yaml")
    want := error(nil)
	if want != err{
		t.Errorf("got %q, wanted %q", err, want)
	}
    defer db.Close()

    err = db.ConnectToDataBase()
    if err != want {
        t.Errorf("got %q, wanted %q", err, want)
    }
	// ЗАМЕНИТЬ
	wantRes:=[]string{"London", "Paris", "Berlin", "New York", "Tokyo"}
	result, errDB := db.GetShortInfoCity("London")
	if errDB != want {
		t.Errorf("got %q, wanted %q", errDB, want)
	}
	if len(wantRes) !=len(result){
		t.Errorf("got %q, wanted %q", result, wantRes)
	}
}

func TestGetFullInfoCity(t *testing.T) {
	db, err := SetConfigs("config/config.yaml")
    want := error(nil)
	if want != err{
		t.Errorf("got %q, wanted %q", err, want)
	}
    defer db.Close()

    err = db.ConnectToDataBase()
    if err != want {
        t.Errorf("got %q, wanted %q", err, want)
    }
	// ЗАМЕНИТЬ
	// wantRes:=[]string{"London", "Paris", "Berlin", "New York", "Tokyo"}
	// result, errDB := db.GetFullInfoCity("London", 2024-06-08)
	// if errDB != want {
	// 	t.Errorf("got %q, wanted %q", errDB, want)
	// }
	// if len(wantRes) !=len(result){
	// 	t.Errorf("got %q, wanted %q", result, wantRes)
	// }
}

