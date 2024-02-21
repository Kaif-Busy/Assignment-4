package Models

import (
	"fmt"

	"github.com/go-pg/pg/v10"
)



func (b *Branch) SaveBranch(db *pg.DB) error {
	_, insertErr := db.Model(b).Returning("*").Insert()
	if insertErr != nil {
		fmt.Println("Error while inserting new item into DB, Reason: &v\n", insertErr)
		return insertErr
	}

	return nil
}
func ShowBranch(db *pg.DB, id uint) (Branch, error) {
	var branch Branch
	err := db.Model(&branch).Where("id=?0", id).Select()
	if err != nil {
		fmt.Println("Error in Select")
		return branch, err
	}
	return branch, nil
}

func ShowAllBranch(db *pg.DB) ([]Branch, error) {
	var branchs []Branch
	err := db.Model(&branchs).Select()
	if err != nil {
		fmt.Println("Error in Select")
		return branchs, err
	}
	// fmt.Printf("Successfully ran Select statement for banks %v\n", banks)
	return branchs, nil
}

func UpdateBranch(db *pg.DB, b *Branch) error {
	_, err := db.Model(b).Where("id=?0", b.ID).UpdateNotZero()
	if err != nil {
		fmt.Println("Error in Update")
		return err
	}
	fmt.Println("Updated.")
	return nil
}
func DeleteBranch(db *pg.DB, b *Branch) error {
	_, err := db.Model(b).Where("id=?0", b.ID).Delete()
	if err != nil {
		fmt.Println("Error in Delete")
		return err
	}
	fmt.Println("Delted.")
	return nil
}
