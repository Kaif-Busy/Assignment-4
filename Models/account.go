package Models

import (
	"errors"
	"fmt"

	"github.com/go-pg/pg/v10"
)

type Account struct {
	ID           uint `pg:",pk"`
	BranchID     uint `pg:",fk:branches,on_delete:SET NULL"`
	AccNo        string
	Balance      float64
	AccType      string
	Branch       *Branch        `pg:"rel:has-one"`
	Transactions []*Transaction `pg:"rel:has-many,on_delete:CASCADE"`
	Customer     []*Customer    `pg:"many2many:customer_to_accounts"`
}

func (b *Account) SaveAcc(db *pg.DB) error {
	_, insertErr := db.Model(b).Returning("*").Insert()
	if insertErr != nil {
		fmt.Println("Error while inserting new item into DB, Reason: &v\n", insertErr)
		return insertErr
	}

	return nil
}
func ShowAcc(db *pg.DB, id uint) (Account, error) {
	var account Account
	err := db.Model(&account).Where("id=?0", id).Select()
	if err != nil {
		fmt.Println("Error in Select")
		return account, err
	}
	return account, nil
}

func ShowAccBalance(db *pg.DB, id uint) (float64, error) {
	var account Account
	err := db.Model(&account).Column("balance").Where("id=?0", id).Select()
	if err != nil {
		fmt.Println("Error in Select")
		return account.Balance, err
	}
	return account.Balance, nil
}

func ShowAllAcc(db *pg.DB) ([]Account, error) {
	var accounts []Account
	err := db.Model(&accounts).Select()
	if err != nil {
		fmt.Println("Error in Select")
		return accounts, err
	}
	// fmt.Printf("Successfully ran Select statement for banks %v\n", banks)
	return accounts, nil
}

func UpdateAcc(db *pg.DB, b *Account) error {
	_, err := db.Model(b).Where("id=?0", b.ID).Update()
	if err != nil {
		fmt.Println("Error in Update")
		return err
	}
	fmt.Println("Updated.")
	return nil
}
func DeleteAcc(db *pg.DB, id uint) error {
	_, err := db.Model(&Account{}).Where("id=?0", id).Delete()
	if err != nil {
		fmt.Println("Error in Delete")
		return err
	}
	fmt.Println("Delted.")
	return nil
}

func Deposit(db *pg.DB, accountID uint, amount float64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	account := &Account{ID: accountID}
	err = tx.Model(account).WherePK().Select()
	if err != nil {
		tx.Rollback()
		return err
	}
	account.Balance += amount
	_, updateErr := tx.Model(account).WherePK().Update()
	if updateErr != nil {
		tx.Rollback()
		return err
	}
	var newTrans Transaction
	newTrans.AccountID=account.ID 
	newTrans.Mode="Cash"
	newTrans.Amount=amount
	newTrans.SaveTransaction(db)
	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func Withdraw(db *pg.DB, accountID uint, amount float64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	account := &Account{ID: accountID}
	err = tx.Model(account).WherePK().Select()
	if err != nil {
		tx.Rollback()
		return err
	}
	if account.Balance < amount {
		tx.Rollback()
		return errors.New("insufficient balance")
	}
	account.Balance -= amount
	_, updateErr := tx.Model(account).WherePK().Update()
	if updateErr != nil {
		tx.Rollback()
		return err
	}
	var newTrans Transaction
	newTrans.AccountID=account.ID 
	newTrans.Mode="ATM"
	newTrans.Amount=-1*amount
	newTrans.SaveTransaction(db)
	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}
