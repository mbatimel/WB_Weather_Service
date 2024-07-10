package migrate

import (
	"context"
	"io"
	"os"

	"github.com/mbatimel/WB_Weather_Service/internal/repo"
)

func ApplyMigrations(db *repo.DataBase, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	queries, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	if !checkingTables(db){
		_, err = db.DB.Exec(context.Background(), string(queries))
		if err != nil {
			return err
		}
	}

	return nil
}

func checkingTables(db *repo.DataBase) bool{
		query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name IN ('cities', 'weather_forecasts', 'persons', 'favorite_cities')
		)
	`

	var exists bool
	err := db.DB.QueryRow(context.Background(), query).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}