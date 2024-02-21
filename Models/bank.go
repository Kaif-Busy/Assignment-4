package Models

import (
	"fmt"

	"github.com/go-pg/pg/v10"
)



func (b *Bank) SaveBank(db *pg.DB) error {
	_, insertErr := db.Model(b).Returning("*").Insert()
	if insertErr != nil {
		fmt.Println("Error while inserting new item into DB, Reason: &v\n", insertErr)
		return insertErr
	}

	return nil
}
func ShowBank(db *pg.DB, id uint) (Bank, error) {
	var bank Bank
	var branches []*Branch
	err := db.Model(&bank).Where("id=?0", id).Select()
	if err != nil {
		fmt.Println("Error in Select")
		return bank, err
	}
	err= db.Model(&branches).Where("bank_id=?0",id).Select()
	if err != nil {
		fmt.Println("Error in Getting Branches for the bank")
		return bank, err
	}
	bank.Branches=branches
	return bank, nil
}

func ShowAllBank(db *pg.DB) ([]Bank, error) {
	var banks []Bank
	err := db.Model(&banks).Select()
	if err != nil {
		fmt.Println("Error in Select")
		return banks, err
	}
	// fmt.Printf("Successfully ran Select statement for banks %v\n", banks)
	return banks, nil
}

func UpdateBank(db *pg.DB, b *Bank) error {
	_, err := db.Model(b).Where("id=?0", b.ID).UpdateNotZero()
	if err != nil {
		fmt.Println("Error in Update")
		return err
	}
	fmt.Println("Updated.")
	return nil
}
func DeleteBank(db *pg.DB, b *Bank) error {
	_, err := db.Model(b).Where("id=?0", b.ID).Delete()
	if err != nil {
		fmt.Println("Error in Delete")
		return err
	}
	fmt.Println("Delted.")
	return nil
}
