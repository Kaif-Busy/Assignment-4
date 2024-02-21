package operations

import (
	"bank/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

type NewAccount struct {
	Acc    Models.Account
	CustID uint
}

func AddAccount(ctx *gin.Context, db *pg.DB) {
	var newAccount NewAccount
	tx, err := db.Begin()
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := ctx.ShouldBindJSON(&newAccount); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := newAccount.Acc.SaveAccount(tx); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to save Account to database"})
		return
	}
	var newCustToAcc Models.CustomerToAccount
	newCustToAcc.CustomerID = newAccount.CustID
	newCustToAcc.AccountID = newAccount.Acc.ID

	if err := newCustToAcc.SaveCustomerToAccount(tx); err != nil {
		Models.DeleteAccount(tx, newAccount.Acc.ID)
		ctx.JSON(500, gin.H{"error": "Failed to save account to database", "details": err.Error()})
		return
	}
	commitErr := tx.Close()
	if commitErr != nil {
		ctx.JSON(500, gin.H{"error": "Failed to close Transaction", "details": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Account added successfully", "data": newAccount})
}

func ViewAllAccounts(ctx *gin.Context, db *pg.DB) {
	accounts, err := Models.ShowAllAccount(db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in printing all accounts",
		})
	} else {
		ctx.JSON(200, accounts)
	}
}

func ViewAccount(ctx *gin.Context, db *pg.DB) {
	naccid, _ := strconv.ParseUint(ctx.Param("accid"), 10, 0)
	account, err := Models.ShowAccount(db, uint(naccid))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in printing all accounts",
		})
	} else {
		ctx.JSON(200, account)
	}
}

func UpdateAccount(ctx *gin.Context, db *pg.DB) {
	var newAccount Models.Account
	if err := ctx.ShouldBindJSON(&newAccount); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := Models.UpdateAccount(db, &newAccount)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(200, newAccount)
	}
}

func ViewBalance(ctx *gin.Context, db *pg.DB) {
	naccid, _ := strconv.ParseUint(ctx.Param("accid"), 10, 0)
	accountBalance, err := Models.ShowAccountBalance(db, uint(naccid))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"balance": accountBalance,
		})
	}

}

func DeleteAccount(ctx *gin.Context, db *pg.DB) {
	var newAccount Models.Account
	tx, txErr := db.Begin()
	if txErr != nil {
		ctx.JSON(400, gin.H{"error": txErr.Error()})
		return
	}
	if err := ctx.ShouldBindJSON(&newAccount); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := Models.DeleteAccount(tx, newAccount.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else {
		commitErr := tx.Close()
		if commitErr != nil {
			ctx.JSON(500, gin.H{"error": "Failed to close Transaction", "details": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{
			"message": "deleted",
		})
	}

}

func DepositMoney(ctx *gin.Context, db *pg.DB) {
	tx, _ := db.Begin()

	naccid, _ := strconv.ParseUint(ctx.Param("accid"), 10, 0)
	namount, _ := strconv.ParseInt(ctx.Param("amount"), 10, 0)

	err := Models.Deposit(tx, uint(naccid), float64(namount))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	commitErr := tx.Commit()
	if commitErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": commitErr.Error(),
		})
	}
	ctx.JSON(200, gin.H{
		"message": "deposit successfull",
	})

}
func WithdrawMoney(ctx *gin.Context, db *pg.DB) {
	tx, _ := db.Begin()
	naccid, _ := strconv.ParseUint(ctx.Param("accid"), 10, 0)
	namount, _ := strconv.ParseInt(ctx.Param("amount"), 10, 0)

	err := Models.Withdraw(tx, uint(naccid), float64(namount))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	commitErr := tx.Commit()
	if commitErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": commitErr.Error(),
		})
	}
	ctx.JSON(200, gin.H{
		"message": "withdraw successfull",
	})

}
