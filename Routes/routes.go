package Routes

import (
	operations "bank/Operations"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

func Routes(db *pg.DB) {
	server := gin.Default()

	//Bank APIs

	server.GET("/bank", func(ctx *gin.Context) { //done
		operations.ViewAllBanks(ctx, db)
	})

	server.GET("/bank/:bid", func(ctx *gin.Context) { //done
		operations.ViewBank(ctx, db)
	})

	server.POST("/bank/", func(ctx *gin.Context) { //done
		operations.AddBank(ctx, db)
	})
	server.PUT("/bank/", func(ctx *gin.Context) { //done
		operations.UpdateBank(ctx, db)
	})

	server.DELETE("/bank/:bid", func(ctx *gin.Context) { //done
		operations.DeleteBank(ctx, db)
	})

	//Branch APIs

	server.GET("/branch", func(ctx *gin.Context) { //done
		operations.ViewAllBranches(ctx, db)
	})

	server.GET("/branch/:bid", func(ctx *gin.Context) { //done
		operations.ViewBranch(ctx, db)
	})

	server.POST("/branch/", func(ctx *gin.Context) { //done
		operations.AddBranch(ctx, db)
	})
	server.PUT("/branch/", func(ctx *gin.Context) { //done
		operations.UpdateBranch(ctx, db)
	})

	server.DELETE("/branch/:bid", func(ctx *gin.Context) { //done
		operations.DeleteBranch(ctx, db)
	})

	//Customer APIs

	server.GET("/customer", func(ctx *gin.Context) {
		operations.ViewAllCustomers(ctx, db)
	})
	server.GET("/customer/:cust_id", func(ctx *gin.Context) {
		operations.ViewCustomer(ctx, db)
	})
	server.GET("/customer/:cust_id/accounts", func(ctx *gin.Context) {
		operations.ViewCustomerAccount(ctx, db)
	})
	server.POST("/customer/", func(ctx *gin.Context) {
		operations.AddCustomer(ctx, db)
	})
	server.POST("/customer/joint_account", func(ctx *gin.Context) {
		operations.AddJointCustomerAccount(ctx, db)
	})
	server.PUT("/customer/:cust_id", func(ctx *gin.Context) {
		operations.UpdateCustomer(ctx, db)
	})
	server.DELETE("/customer/:cust_id", func(ctx *gin.Context) {
		operations.DeleteCustomer(ctx, db)
	})

	//Account APIs

	server.GET("/account", func(ctx *gin.Context) {
		operations.ViewAllAccounts(ctx, db)
	})
	server.GET("/account/:accid", func(ctx *gin.Context) {
		operations.ViewAccount(ctx, db)
	})
	server.GET("/account/balance/:accid", func(ctx *gin.Context) {
		operations.ViewBalance(ctx, db)
	})
	server.POST("/account/", func(ctx *gin.Context) {
		operations.AddAccount(ctx, db)
	})
	server.PUT("/account/:accid", func(ctx *gin.Context) {
		operations.UpdateAccount(ctx, db)
	})
	server.POST("/account/:accid/deposit/:amount", func(ctx *gin.Context) {
		operations.DepositMoney(ctx, db)
	})
	server.POST("/account/:accid/withdraw/:amount", func(ctx *gin.Context) {
		operations.WithdrawMoney(ctx, db)
	})
	server.DELETE("/account/:accid", func(ctx *gin.Context) {
		operations.DeleteAccount(ctx, db)
	})

	// //trasnfer Money
	server.POST("/transaction/transfer", func(ctx *gin.Context) {
		operations.AddTransaction(ctx, db)
	})
	server.Run()
}
