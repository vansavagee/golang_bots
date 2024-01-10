package sqlqueries

import (
	"database/sql"
	"listen_bot/internal/model"
	"os"
	"time"
)

func InsertData(v model.Device) error {
	db, err := sql.Open("postgres", os.Getenv("ConnStr"))
	if err != nil {
		return err
	}
	// Пример операции INSERT
	stmt, err := db.Prepare(" INSERT INTO storage (company, model, producer, country, price, date) VALUES ($1, $2,$3, $4,$5,$6);")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(v.Company, v.Model, v.Customer, v.Country, v.Price, time.Now().Local().Format("2006-01-02"))
	if err != nil {
		return err
	}
	return nil
}
