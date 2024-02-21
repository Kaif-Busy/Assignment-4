package operations

import (
	"bank/Models"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

func AddTransaction(ctx *gin.Context, db *pg.DB) {
	var newTransaction Models.Transaction
	tx,_:=db.Begin()
	if err := ctx.ShouldBindJSON(&newTransaction); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := newTransaction.Trans(tx); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	commitErr := tx.Commit()
	if commitErr != nil {
		ctx.JSON(500, gin.H{"error": commitErr.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Account added successfully", "data": newTransaction})
}
