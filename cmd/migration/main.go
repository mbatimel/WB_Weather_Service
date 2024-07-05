package main

import (
	"log"

	"github.com/mbatimel/WB_Weather_Service/internal/migrate"
	"github.com/mbatimel/WB_Weather_Service/internal/repo"
)

func main(){
	db, err := repo.SetConfigs("config/config.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	db.ConnectToDataBase()
	defer db.Close()
	if err = migrate.CreateSchema(db); err != nil {
		log.Fatalln(err)
	}

}