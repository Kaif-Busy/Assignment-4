package Models

import "time"

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

type Bank struct {
	ID       uint `pg:",pk"`
	BankName string
	Branches []*Branch `pg:"rel:has-many,on_delete:CASCADE"`
}

type Branch struct {
	ID        uint `pg:",pk"`
	Address   string
	IFSCCode  string
	BankID    uint        `pg:",fk:bank_id,on_delete:SET NULL"`
	Bank      Bank        `pg:"rel:has-one"`
	Customers []*Customer `pg:"rel:has-many,on_delete:CASCADE"`
}

type Customer struct {
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

type CustomerToAccount struct {
	ID	uint
	CustomerID uint      `pg:"on_delete:CASCADE"`
	Customer   *Customer `pg:"rel:has-one"`
	AccountID  uint      `pg:"on_delete:CASCADE"`
	Account    *Account  `pg:"rel:has-one"`
}

type Transaction struct {
	ID                uint `pg:",pk"`
	AccountID         uint `pg:",fk:accounts,on_delete:SET NULL"`
	Mode              string
	ReceiverAccNumber string
	Timestamp         time.Time `pg:",default:now()"`
	Amount            float64
	Account           *Account `pg:"rel:has-one"`
}
