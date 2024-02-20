package main

import (
	"bank/Models"
	"bank/Routes"
	"fmt"
	_ "net/http"
	"os"

	_ "github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

func Connect() *pg.DB {
	opts := &pg.Options{
		User:     "kaif",
		Password: "kaif",
		Database: "mydb",
		Addr:     "localhost:5432",
	}

	db := pg.Connect(opts)
	if db == nil {
		fmt.Print("failed to connect to the database")
		os.Exit(100)
	} else {

		fmt.Println("Hello, you are connected to the database")
	}

	err := CreateBankTables(db)
	if err != nil {
		panic(err.Error())
	}

	return db
}

func init() {
	orm.RegisterTable((*Models.CustomerToAccount)(nil))
}

func CreateBankTables(db *pg.DB) error {
	models := []interface{}{
		(*Models.Bank)(nil),
		(*Models.Branch)(nil),
		(*Models.Customer)(nil),
		(*Models.Account)(nil),
		(*Models.CustomerToAccount)(nil),
		(*Models.Transaction)(nil),
	}
	opts := &orm.CreateTableOptions{
		IfNotExists:   true,
		FKConstraints: true,
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(opts)
		if err != nil {
			return err
		}

	}

	return nil

}

func main() {
	db:=Connect()
	Routes.Routes(db)
}
