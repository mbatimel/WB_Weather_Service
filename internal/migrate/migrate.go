package migrate

import 
(
	"github.com/go-pg/pg/orm"
	"github.com/mbatimel/WB_Weather_Service/internal/repo"
	"github.com/mbatimel/WB_Weather_Service/internal/model"

)


func CreateSchema(db *repo.DataBase) error {
	models := []interface{}{
		(*model.WeatherForecast)(nil),
		(*model.Cities)(nil),
	}
	for _, model := range models {
		op := orm.CreateTableOptions{}
		err := db.DB.Model(model).CreateTable(&op)
		if err != nil {
			return err
		}
	}
	return nil
}