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

var productsFilePath = "docs/db/json/products.json"

// NewProductsLoader loads products data
func NewProductsLoader(db *sql.DB) *ProductLoader {
	return &ProductLoader{path: productsFilePath, db: db}
}

// ProductLoader is a struct that returns the invoice loader
type ProductLoader struct {
	path string
	db   *sql.DB
}

func (cl *ProductLoader) ReadAll() (err error) {
	jsonFile, err := os.Open(cl.path)

	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		log.Fatal(err)
	}

	var products []handler.ProductJSON
	err = json.Unmarshal(byteValue, &products)
	if err != nil {
		return err
	}

	// Insert in database
	for _, product := range products {
		_, err := cl.db.Exec(
			"INSERT INTO `products` (`id`, `description`, `price`) VALUES (?, ?, ?)",
			product.Id, product.Description, product.Price,
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
