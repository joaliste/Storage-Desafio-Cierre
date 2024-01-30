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

var salesFilePath = "docs/db/json/sales.json"

func NewSalesLoader(db *sql.DB) *SalesLoader {
	return &SalesLoader{path: salesFilePath, db: db}
}

type SalesLoader struct {
	path string
	db   *sql.DB
}

func (cl *SalesLoader) ReadAll() (err error) {
	jsonFile, err := os.Open(cl.path)

	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		log.Fatal(err)
	}

	var sales []handler.SaleJSON
	err = json.Unmarshal(byteValue, &sales)
	if err != nil {
		return err
	}

	// Insert in database
	for _, sale := range sales {
		_, err := cl.db.Exec(
			"INSERT INTO `sales` (`id`, `quantity`, `invoice_id`, `product_id`) VALUES (?, ?, ?, ?)",
			sale.Id, sale.Quantity, sale.InvoiceId, sale.ProductId,
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
