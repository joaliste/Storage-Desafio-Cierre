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

var invoicesFilePath = "docs/db/json/invoices.json"

// NewInvoiceLoader loads invoice data
func NewInvoiceLoader(db *sql.DB) *InvoiceLoader {
	return &InvoiceLoader{path: invoicesFilePath, db: db}
}

// InvoiceLoader is a struct that returns the invoice loader
type InvoiceLoader struct {
	path string
	db   *sql.DB
}

func (cl *InvoiceLoader) ReadAll() (err error) {
	jsonFile, err := os.Open(cl.path)

	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		log.Fatal(err)
	}

	var invoices []handler.InvoiceJSON
	err = json.Unmarshal(byteValue, &invoices)
	if err != nil {
		return err
	}

	// Insert in database
	for _, invoice := range invoices {
		_, err := cl.db.Exec(
			"INSERT INTO `invoices` (`id`, `datetime`, `customer_id`, `total`) VALUES (?, ?, ?, ?)",
			invoice.Id, invoice.Datetime, invoice.CustomerId, invoice.Total,
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
