package Models

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
)

type Transaction struct {
	// tableName         struct{}  `pg:"transactions"`
	ID                uint `pg:",pk"`
	AccountID         uint `pg:",fk:accounts,on_delete:SET NULL"`
	Mode              string
	ReceiverAccNumber string
	Timestamp         time.Time `pg:",default:now()"`
	Amount            float64
	Account           *Account `pg:"rel:has-one"`
}

func (tr *Transaction) SaveTransaction(db *pg.DB) error {
	_, insertErr := db.Model(tr).Returning("*").Insert()
	if insertErr != nil {
		fmt.Println("Error while inserting new item into DB, Reason: &v\n", insertErr)
		return insertErr
	}

	return nil
}

func (tr *Transaction) Trans(db *pg.DB) error {
	senderID := tr.AccountID
	amount := tr.Amount
	receiverAccountNo := tr.ReceiverAccNumber

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Fetch sender's account
	senderAccount := &Account{ID: senderID}
	err = tx.Model(senderAccount).WherePK().Select()
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
	tr.SaveTransaction(db)
	// Commit the transaction
	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil

}
