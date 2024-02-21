package Models

import (
	"errors"
	"fmt"

	"github.com/go-pg/pg/v10"
)

func (b *Account) SaveAccount(tx *pg.Tx) error {
	_, insertErr := tx.Model(b).Returning("*").Insert()
	if insertErr != nil {
		fmt.Println("Error while inserting new item into DB, Reason: &v\n", insertErr)
		tx.Rollback()
		return insertErr
	}

	return nil
}
func ShowAccount(db *pg.DB, id uint) (Account, error) {
	var account Account
	err := db.Model(&account).Where("id=?0", id).Select()
	if err != nil {
		fmt.Println("Error in Select")
		return account, err
	}
	return account, nil
}

func ShowAccountBalance(db *pg.DB, id uint) (float64, error) {
	var account Account
	err := db.Model(&account).Column("balance").Where("id=?0", id).Select()
	if err != nil {
		fmt.Println("Error in Select")
		return account.Balance, err
	}
	return account.Balance, nil
}

func ShowAllAccount(db *pg.DB) ([]Account, error) {
	var accounts []Account
	err := db.Model(&accounts).Select()
	if err != nil {
		fmt.Println("Error in Select")
		return accounts, err
	}
	return accounts, nil
}

func UpdateAccount(db *pg.DB, b *Account) error {
	_, err := db.Model(b).Where("id=?0", b.ID).UpdateNotZero()
	if err != nil {
		fmt.Println("Error in Update")
		return err
	}
	fmt.Println("Updated.")
	return nil
}
func DeleteAccount(tx *pg.Tx, id uint) error {
	if _, err := tx.Model(&Account{}).Where("id=?0", id).Delete(); err != nil {
		fmt.Println("Error in Delete")
		tx.Rollback()
		return err
	}
	if _, err := tx.Model(&CustomerToAccount{}).Where("account_id=?0", id).Delete(); err != nil {
		fmt.Println("Error in Delete")
		tx.Rollback()
		return err
	}
	fmt.Println("Delted.")
	return nil
}

func Deposit(tx *pg.Tx, accountID uint, amount float64) error {

	account := &Account{ID: accountID}
	err := tx.Model(account).WherePK().Select()
	if err != nil {
		tx.Rollback()
		return err
	}
	account.Balance += amount
	_, updateErr := tx.Model(account).WherePK().Update()
	if updateErr != nil {
		tx.Rollback()
		return updateErr
	}
	var newTrans Transaction
	newTrans.AccountID = account.ID
	newTrans.Mode = "Cash"
	newTrans.Amount = amount
	newTrans.SaveTransaction(tx)
	return nil
}

func Withdraw(tx *pg.Tx, accountID uint, amount float64) error {

	account := &Account{ID: accountID}
	err := tx.Model(account).WherePK().Select()
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
		return updateErr
	}
	var newTrans Transaction
	newTrans.AccountID = account.ID
	newTrans.Mode = "ATM"
	newTrans.Amount = -1 * amount
	newTrans.SaveTransaction(tx)

	return nil
}
