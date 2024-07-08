package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/mbatimel/WB_Weather_Service/internal/config"
)

type DataBase struct {
	DB     *pgx.Conn
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
	if db.DB != nil {
		db.DB.Close(context.Background())
	}
}

func (db *DataBase) ConnectToDataBase() error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db.config.User, db.config.Password, db.config.Host, db.config.Port, db.config.Database)
    conn, err := pgx.Connect(context.Background(), connStr)
    if err!=nil{
        return err
    }
	db.DB = conn
    return nil
}
