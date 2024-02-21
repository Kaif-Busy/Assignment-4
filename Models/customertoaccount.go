package Models

import (
	"fmt"

	"github.com/go-pg/pg/v10"
)



func (cta *CustomerToAccount) SaveCustomerToAccount(tx *pg.Tx) error {
	_, insertErr := tx.Model(cta).Returning("*").Insert()
	if insertErr != nil {
		fmt.Println("Error while inserting new relation between customer and account into DB, Reason: &v\n", insertErr)
		tx.Rollback()
		return insertErr
	}
	return nil
}
