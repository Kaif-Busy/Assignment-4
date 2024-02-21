package Models

import (
	"errors"
	"fmt"

	"github.com/go-pg/pg/v10"
)

func (tr *Transaction) SaveTransaction(tx *pg.Tx) error {
	_, insertErr := tx.Model(tr).Returning("*").Insert()
	if insertErr != nil {
		fmt.Println("Error while inserting new item into DB, Reason: &v\n", insertErr)
		return insertErr
	}

	return nil
}

func (tr *Transaction) Trans(tx *pg.Tx) error {
	senderID := tr.AccountID
	amount := tr.Amount
	receiverAccountNo := tr.ReceiverAccNumber

	// Fetch sender's account
	senderAccount := &Account{ID: senderID}
	err := tx.Model(senderAccount).WherePK().Select()
	if err != nil {
		tx.Rollback()
		return errors.New("sender account issue")
	}

	// Check if sender has sufficient balance
	if senderAccount.Balance < amount {
		tx.Rollback()
		return errors.New("insufficient balance")
	}

	// Deduct amount from sender's account
	senderAccount.Balance -= amount
	_, updateErr := tx.Model(senderAccount).WherePK().Update()
	if updateErr != nil {
		tx.Rollback()
		return updateErr
	}

	// Fetch receiver's account
	var receiverAccount Account
	err = tx.Model(&receiverAccount).Where("acc_no = ?0", receiverAccountNo).Select()
	if err != nil {
		tx.Rollback()
		return errors.New("receiver account issue")
	}

	// Credit amount to receiver's account
	receiverAccount.Balance += amount
	_, updateErr = tx.Model(&receiverAccount).WherePK().Update()
	if updateErr != nil {
		tx.Rollback()
		return updateErr
	}
	tr.SaveTransaction(tx)

	return nil

}
