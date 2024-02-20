package Models

import (
	"fmt"

	"github.com/go-pg/pg/v10"
)

type Bank struct {
	//tableName struct{}  `pg:"banks"`
	ID       uint `pg:",pk"`
	BankName string
	Branches []*Branch `pg:"rel:has-many,on_delete:CASCADE"`
}

func (b *Bank) SaveB(db *pg.DB) error {
	_, insertErr := db.Model(b).Returning("*").Insert()
	if insertErr != nil {
		fmt.Println("Error while inserting new item into DB, Reason: &v\n", insertErr)
		return insertErr
	}

	return nil
}
func ShowB(db *pg.DB, id uint) (Bank, error) {
	var bank Bank
	err := db.Model(&bank).Where("id=?0", id).Select()
	if err != nil {
		fmt.Println("Error in Select")
		return bank, err
	}
	return bank, nil
}

func ShowAllB(db *pg.DB) ([]Bank, error) {
	var banks []Bank
	err := db.Model(&banks).Select()
	if err != nil {
		fmt.Println("Error in Select")
		return banks, err
	}
	// fmt.Printf("Successfully ran Select statement for banks %v\n", banks)
	return banks, nil
}

func UpdateB(db *pg.DB, b *Bank) error {
	_, err := db.Model(b).Where("id=?0", b.ID).Update()
	if err != nil {
		fmt.Println("Error in Update")
		return err
	}
	fmt.Println("Updated.")
	return nil
}
func DeleteB(db *pg.DB, b *Bank) error {
	_, err := db.Model(b).Where("id=?0", b.ID).Delete()
	if err != nil {
		fmt.Println("Error in Delete")
		return err
	}
	fmt.Println("Delted.")
	return nil
}
