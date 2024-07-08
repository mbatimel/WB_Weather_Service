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

	_, err = db.DB.Exec(context.Background(), string(queries))
	if err != nil {
		return err
	}

	return nil
}
