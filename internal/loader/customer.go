package loader

import (
	"app/internal/handler"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-sql-driver/mysql"
	"io"
	"log"
	"os"
)

var filePath = "docs/db/json/customers.json"

// NewCustomerLoader loads customer data
func NewCustomerLoader(db *sql.DB) *CustomersLoader {
	return &CustomersLoader{path: filePath, db: db}
}

// CustomersLoader is a struct that returns the customer handlers
type CustomersLoader struct {
	path string
	db   *sql.DB
}

func (cl *CustomersLoader) ReadAll() (err error) {
	jsonFile, err := os.Open(cl.path)

	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		log.Fatal(err)
	}

	var customers []handler.CustomerJSON
	err = json.Unmarshal(byteValue, &customers)
	if err != nil {
		return err
	}

	// Insert in database
	for _, customer := range customers {
		_, err := cl.db.Exec(
			"INSERT INTO `customers` (`id`, `first_name`, `last_name`, `condition`) VALUES (?, ?, ?, ?)",
			customer.Id, customer.FirstName, customer.LastName, customer.Condition,
		)

		if err != nil {
			log.Println(err.Error())
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) {
				if mysqlErr.Number != 1062 {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
