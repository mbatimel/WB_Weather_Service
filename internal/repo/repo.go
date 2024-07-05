package repo

import (
	"github.com/go-pg/pg"
	"github.com/mbatimel/WB_Weather_Service/internal/config"
)

type DataBase struct {
	DB     *pg.DB
	config *config.Repo
}

func SetConfigs(path string) (*DataBase, error) {
	config, err := config.NewConfigDB(path)
	if err != nil {
		return nil, err
	}
	return &DataBase{nil, config}, nil
}

func (db *DataBase) Close() {
	db.DB.Close()
}

func (db *DataBase) ConnectToDataBase() {
	db.DB = pg.Connect(
		&pg.Options{
			User:     db.config.User,
			Password: db.config.Password,
			Database: db.config.Database,
		})
}
