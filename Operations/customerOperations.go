package operations

import (
	"bank/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

type NewCustomer struct {
	Acc Models.Account
	Cu  Models.Customer
}

type JointCustomer struct {
	Cu    Models.Customer
	AccID uint
}

func AddCustomer(ctx *gin.Context, db *pg.DB) {
	var newCustomer NewCustomer
	if err := ctx.ShouldBindJSON(&newCustomer); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid JSON format or missing required fields"})
		return
	}

	// Save customer
	if err := newCustomer.Cu.SaveC(db); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to save customer to database", "details": err.Error()})
		return
	}

	// Save account
	if err := newCustomer.Acc.SaveAcc(db); err != nil {
		// If there's an error saving the account, consider rolling back the customer creation to maintain consistency.
		Models.DeleteC(db, &newCustomer.Cu)
		ctx.JSON(500, gin.H{"error": "Failed to save account to database", "details": err.Error()})
		return
	}

	var newCustToAcc Models.CustomerToAccount
	newCustToAcc.CustomerID = newCustomer.Cu.ID
	newCustToAcc.AccountID = newCustomer.Acc.ID

	if err := newCustToAcc.SaveCusToAcc(db); err != nil {
		Models.DeleteC(db, &newCustomer.Cu)
		Models.DeleteAcc(db, newCustomer.Acc.ID)
		ctx.JSON(500, gin.H{"error": "Failed to save account to database", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Customer and account added successfully", "data": newCustomer})

}

func AddJointCustomerAccount(ctx *gin.Context, db *pg.DB) {
	var newCustomer JointCustomer
	if err := ctx.ShouldBindJSON(&newCustomer); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid JSON format or missing required fields"})
		return
	}
	if err := newCustomer.Cu.SaveC(db); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to save customer to database", "details": err.Error()})
		return
	}
	var newCustToAcc Models.CustomerToAccount
	newCustToAcc.CustomerID = newCustomer.Cu.ID
	newCustToAcc.AccountID = newCustomer.AccID

	if err := newCustToAcc.SaveCusToAcc(db); err != nil {
		Models.DeleteC(db, &newCustomer.Cu)
		ctx.JSON(500, gin.H{"error": "Failed to save account to database", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Customer and account added successfully", "data": newCustomer})
}

func ViewAllCustomers(ctx *gin.Context, db *pg.DB) {
	Customers, err := Models.ShowAllC(db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in printing all Customers",
		})
	} else {
		ctx.JSON(200, Customers)
	}
}

func ViewCustomer(ctx *gin.Context, db *pg.DB) {
	nbid, _ := strconv.ParseUint(ctx.Param("cust_id"), 10, 0)
	customer, err := Models.ShowC(db, uint(nbid))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in printing all Customers",
		})
	} else {
		ctx.JSON(200, customer)
	}
}

func ViewCustomerAccount(ctx *gin.Context, db *pg.DB) {
	nbid, _ := strconv.ParseUint(ctx.Param("cust_id"), 10, 0)
	accounts, err := Models.ShowCustomerBanks(db, uint(nbid))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, accounts)
	}
}

func UpdateCustomer(ctx *gin.Context, db *pg.DB) {
	var newCustomer Models.Customer
	if err := ctx.ShouldBindJSON(&newCustomer); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := Models.UpdateC(db, &newCustomer)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(200, newCustomer)
	}
}

func DeleteCustomer(ctx *gin.Context, db *pg.DB) {
	var newCustomer Models.Customer
	if err := ctx.ShouldBindJSON(&newCustomer); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := Models.DeleteC(db, &newCustomer)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in deleting Customers",
		})
	} else {
		ctx.JSON(200, gin.H{
			"message": "deleted",
		})
	}
}
