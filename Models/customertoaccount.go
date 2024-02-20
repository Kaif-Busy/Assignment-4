package Models

import (
	"fmt"

	"github.com/go-pg/pg/v10"
)

type CustomerToAccount struct {
	CustomerID uint
	Customer   *Customer `pg:"rel:has-one"`
	AccountID  uint
	Account    *Account `pg:"rel:has-one"`
}

func (cta *CustomerToAccount) SaveCusToAcc(db *pg.DB) error {
	_, insertErr := db.Model(cta).Returning("*").Insert()
	if insertErr != nil {
		fmt.Println("Error while inserting new relation between customer and account into DB, Reason: &v\n", insertErr)
		return insertErr
	}
	return nil
}
