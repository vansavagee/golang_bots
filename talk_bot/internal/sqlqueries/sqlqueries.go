package sqlqueries

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func SelectUniqCompanies() ([]string, error) {
	v := make([]string, 0)
	db, err := sql.Open("postgres", os.Getenv("ConnStr"))
	if err != nil {
		return v, err
	}
	defer db.Close()
	// Пример операции INSERT
	rows, err := db.Query("SELECT company FROM storage GROUP BY company;")
	if err != nil {
		return v, err
	}
	defer rows.Close()
	for rows.Next() {
		var row string
		err = rows.Scan(&row)
		v = append(v, row)
		if err != nil {
			return v, err
		}

	}

	err = rows.Err()
	if err != nil {
		return v, err
	}
	return v, nil
}
func DeleteData() error {
	query := `DELETE FROM storage
	WHERE date <= CURRENT_DATE - INTERVAL '2 day';`
	db, err := sql.Open("postgres", os.Getenv("ConnStr"))
	if err != nil {
		return err
	}
	defer db.Close()
	db.Query(query)
	return nil
}
func SelectAll() (string, error) {
	v := ""
	db, err := sql.Open("postgres", os.Getenv("ConnStr"))
	if err != nil {
		return v, err
	}
	defer db.Close()

	query := "SELECT id,company,model FROM storage;"
	rows, err := db.Query(query)
	if err != nil {
		return v, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id, company, model string
		)
		err = rows.Scan(&id, &company, &model)
		v += fmt.Sprintf("%v  %v %v\n", id, company, model)
		if err != nil {
			return v, err
		}

	}

	err = rows.Err()
	if err != nil {
		return v, err
	}
	return v, nil
}
func SelectAllByIdsArray(ids []int) (string, error) {
	v := ""
	db, err := sql.Open("postgres", os.Getenv("ConnStr"))
	if err != nil {
		return v, err
	}
	defer db.Close()
	for _, id := range ids {
		query := "SELECT company,model,producer,country,price,date FROM storage WHERE id = $1 ORDER BY price;"
		rows, err := db.Query(query, id)
		if err != nil {
			return v, err
		}
		defer rows.Close()
		for rows.Next() {
			var (
				company, model, producer, country, price string
				date                                     time.Time
			)
			err = rows.Scan(&company, &model, &producer, &country, &price, &date)
			year, month, day := date.Date()
			v += fmt.Sprintf("\n%v %v %v %v-%v-%v %v %v\n", company, model, country, year, month, day, producer, price)
			if err != nil {
				return v, err
			}

		}

		err = rows.Err()
		if err != nil {
			return v, err
		}
	}

	return v, nil
}
func SelectUniqDevices(company string) ([]string, error) {
	v := make([]string, 0)
	db, err := sql.Open("postgres", os.Getenv("ConnStr"))
	if err != nil {
		return v, err
	}
	defer db.Close()
	query := "SELECT model FROM storage WHERE company = $1 GROUP BY model;"
	rows, err := db.Query(query, company)
	if err != nil {
		return v, err
	}
	defer rows.Close()
	for rows.Next() {
		var row string
		err = rows.Scan(&row)
		v = append(v, row)
		if err != nil {
			return v, err
		}

	}

	err = rows.Err()
	if err != nil {
		return v, err
	}
	return v, nil
}
