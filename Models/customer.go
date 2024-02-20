package Models

import (
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
)

type Customer struct {
	//tableName   struct{}  `pg:"customer"`
	ID          uint `pg:",pk"`
	Name        string
	PANNumber   string
	DOB         time.Time
	Age         int
	PhoneNumber string
	Address     string
	BranchID    uint       `pg:"fk:branch_id,on_delete:SET NULL"`
	Branch      *Branch    `pg:"rel:has-one"`
	Account     []*Account `pg:"many2many:customer_to_accounts,on_delete:CASCADE"`
}

func (b *Customer) SaveC(db *pg.DB) error {
	_, insertErr := db.Model(b).Returning("*").Insert()
	if insertErr != nil {
		fmt.Println("Error while inserting new item into DB, Reason: &v\n", insertErr)
		return insertErr
	}

	return nil
}
func ShowC(db *pg.DB, id uint) (Customer, error) {
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
    return nil,err
}

	return customer.Account, nil
}

func ShowAllC(db *pg.DB) ([]Customer, error) {
	var customers []Customer
	err := db.Model(&customers).Select()
	if err != nil {
		fmt.Println("Error in Select")
		return customers, err
	}
	return customers, nil
}

func UpdateC(db *pg.DB, b *Customer) error {
	_, err := db.Model(b).Where("id=?0", b.ID).Update()
	if err != nil {
		fmt.Println("Error in Update")
		return err
	}
	fmt.Println("Updated.")
	return nil
}
func DeleteC(db *pg.DB, b *Customer) error {
	_, err := db.Model(b).Where("id=?0", b.ID).Delete()
	if err != nil {
		fmt.Println("Error in Delete")
		return err
	}
	fmt.Println("Delted.")
	return nil
}
