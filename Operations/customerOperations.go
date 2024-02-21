package operations

import (
	"bank/Models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

type DisplayCustomer struct {
	ID          uint
	Name        string
	PANNumber   string
	DOB         time.Time
	Age         int
	PhoneNumber string
	Address     string
}

type NewCustomer struct {
	Account  Models.Account
	Customer Models.Customer
}

type JointCustomer struct {
	Customer Models.Customer
	AccID    uint
}

func AddCustomer(ctx *gin.Context, db *pg.DB) {
	var newCustomer NewCustomer
	tx, err := db.Begin()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := ctx.ShouldBindJSON(&newCustomer); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid JSON format or missing required fields"})
		return
	}

	// Save customer
	if err := newCustomer.Customer.SaveCustomer(tx); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to save customer to database", "details": err.Error()})
		return
	}

	// Save account
	if err := newCustomer.Account.SaveAccount(tx); err != nil {
		Models.DeleteCustomer(tx, &newCustomer.Customer)
		ctx.JSON(500, gin.H{"error": "Failed to save account to database", "details": err.Error()})
		return
	}

	var newCustToAcc Models.CustomerToAccount
	newCustToAcc.CustomerID = newCustomer.Customer.ID
	newCustToAcc.AccountID = newCustomer.Account.ID

	if err := newCustToAcc.SaveCustomerToAccount(tx); err != nil {
		Models.DeleteCustomer(tx, &newCustomer.Customer)
		Models.DeleteAccount(tx, newCustomer.Account.ID)
		ctx.JSON(500, gin.H{"error": "Failed to save account to database", "details": err.Error()})
		return
	}
	commitErr := tx.Close()
	if commitErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": commitErr.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{"message": "Customer and account added successfully", "Id": newCustomer.Customer.ID})

}

func AddJointCustomerAccount(ctx *gin.Context, db *pg.DB) {
	var newCustomer JointCustomer
	tx, err := db.Begin()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := ctx.ShouldBindJSON(&newCustomer); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid JSON format or missing required fields"})
		return
	}
	if err := newCustomer.Customer.SaveCustomer(tx); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to save customer to database", "details": err.Error()})
		return
	}
	var newCustToAcc Models.CustomerToAccount
	newCustToAcc.CustomerID = newCustomer.Customer.ID
	newCustToAcc.AccountID = newCustomer.AccID

	if err := newCustToAcc.SaveCustomerToAccount(tx); err != nil {
		Models.DeleteCustomer(tx, &newCustomer.Customer)
		ctx.JSON(500, gin.H{"error": "Failed to save account to database", "details": err.Error()})
		return
	}
	commitErr := tx.Close()
	if commitErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": commitErr.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{"message": "Customer and account added successfully", "ID": newCustomer.Customer.ID})
}

func ViewAllCustomers(ctx *gin.Context, db *pg.DB) {
	customers, err := Models.ShowAllCustomer(db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in printing all Customers",
		})
	} else {
		var displayCustomers []DisplayCustomer
		for _, customer := range customers {
			var displayCustomer DisplayCustomer
			displayCustomer.ID = customer.ID
			displayCustomer.Name = customer.Name
			displayCustomer.DOB = customer.DOB
			displayCustomer.Age = customer.Age
			displayCustomer.PANNumber = customer.PANNumber
			displayCustomer.PhoneNumber = customer.PhoneNumber
			displayCustomer.Address = customer.Address

			displayCustomers = append(displayCustomers, displayCustomer)
		}
		ctx.JSON(200, displayCustomers)
	}
}

func ViewCustomer(ctx *gin.Context, db *pg.DB) {
	newCustId, _ := strconv.ParseUint(ctx.Param("cust_id"), 10, 0)
	customer, err := Models.ShowCustomer(db, uint(newCustId))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in printing all Customers",
		})
	} else {
		var displayCustomer DisplayCustomer
		displayCustomer.ID = customer.ID
		displayCustomer.Name = customer.Name
		displayCustomer.DOB = customer.DOB
		displayCustomer.Age = customer.Age
		displayCustomer.PANNumber = customer.PANNumber
		displayCustomer.PhoneNumber = customer.PhoneNumber
		displayCustomer.Address = customer.Address
		ctx.JSON(200, displayCustomer)
	}
}

func ViewCustomerAccount(ctx *gin.Context, db *pg.DB) {
	newCustId, _ := strconv.ParseUint(ctx.Param("cust_id"), 10, 0)
	accounts, err := Models.ShowCustomerBanks(db, uint(newCustId))

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
	err := Models.UpdateCustomer(db, &newCustomer)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(200, gin.H{"message": "Updated"})
	}
}

func DeleteCustomer(ctx *gin.Context, db *pg.DB) {
	var newCustomer Models.Customer
	tx, err := db.Begin()
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := ctx.ShouldBindJSON(&newCustomer); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = Models.DeleteCustomer(tx, &newCustomer)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else {
		commitErr := tx.Close()
		if commitErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": commitErr.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"message": "deleted",
		})
	}
}
