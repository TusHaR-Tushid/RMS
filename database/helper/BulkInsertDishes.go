package helper

import (
	"RMS/database"
	"RMS/models"

	"github.com/elgris/sqrl"
	"github.com/sirupsen/logrus"
)

func BulkInsertDishes(bulkDishes []models.BulkDishes, dishID int) error {
	psql := sqrl.StatementBuilder.PlaceholderFormat(sqrl.Dollar)
	sql := psql.Insert("dish_per_image").Columns("dish_id", "image_id")
	for _, post := range bulkDishes {
		sql.Values(dishID, post.ImageID)
	}

	SQL, args, err := sql.ToSql()
	if err != nil {
		logrus.Printf("BulkInsertDishes: not able to create sql string:%v", err)
		return err
	}
	_, err = database.RmsDB.Exec(SQL, args...)
	if err != nil {
		logrus.Printf("BulkInsertDishes:not able to create bulk dishes:%v", err)
		return err
	}

	return nil
}
