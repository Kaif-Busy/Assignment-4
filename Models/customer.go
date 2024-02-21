package Models

import (
	"fmt"

	"github.com/go-pg/pg/v10"
)

func (b *Customer) SaveCustomer(tx *pg.Tx) error {
	_, insertErr := tx.Model(b).Returning("*").Insert()
	if insertErr != nil {
		fmt.Println("Error while inserting new item into DB, Reason: &v\n", insertErr)
		tx.Rollback()
		return insertErr
	}

	return nil
}
func ShowCustomer(db *pg.DB, id uint) (Customer, error) {
	var customer Customer
	err := db.Model(&customer).Where("id=?0", id).Select()
	if err != nil {
		fmt.Println("Error in Select")
		return customer, err
	}
	return customer, nil
}
func ShowCustomerBanks(db *pg.DB, id uint) ([]*Account, error) {
	var customer Customer
	err := db.Model(&customer).
		Relation("Account").
		Where("customer.id = ?", id).
		Select()

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return customer.Account, nil
}

func ShowAllCustomer(db *pg.DB) ([]Customer, error) {
	var customers []Customer
	err := db.Model(&customers).Select()
	if err != nil {
		fmt.Println("Error in Select")
		return customers, err
	}
	return customers, nil
}

func UpdateCustomer(db *pg.DB, b *Customer) error {
	_, err := db.Model(b).Where("id=?0", b.ID).UpdateNotZero()
	if err != nil {
		fmt.Println("Error in Update")
		return err
	}
	fmt.Println("Updated.")
	return nil
}
func DeleteCustomer(tx *pg.Tx, b *Customer) error {
	var accountsRealted []CustomerToAccount
	tx.Model(&accountsRealted).Where("customer_id=?0", b.ID).Returning("account_id").Select()

	if _, err := tx.Model(b).Where("id=?0", b.ID).Delete(); err != nil {
		fmt.Println("Error in Delete")
		tx.Rollback()
		return err
	}
	if _, err := tx.Model(&CustomerToAccount{}).Where("customer_id=?0", b.ID).Delete(); err != nil {
		fmt.Println("Error in Delete")
		tx.Rollback()
		return err
	}
	for _, accountRealted := range accountsRealted {
		rowsAffected, err := tx.Model(&CustomerToAccount{}).Where("account_id=?0", accountRealted.AccountID).Count()
		if err != nil {
			tx.Rollback()
			return err
		}
		if rowsAffected == 0 {
			DeleteAccount(tx, accountRealted.AccountID)
		}

	}

	fmt.Println("Delted.")
	return nil
}
