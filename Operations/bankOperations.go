package operations

import (
	"bank/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

func AddBank(ctx *gin.Context, db *pg.DB) {
	var newBank Models.Bank
	if err := ctx.ShouldBindJSON(&newBank); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := newBank.SaveB(db); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to save Bank to database"})
		return
	}
	ctx.JSON(200, gin.H{"message": "Bank added successfully", "data": newBank})
}

func ViewAllBanks(ctx *gin.Context, db *pg.DB) {
	banks, err := Models.ShowAllB(db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in printing all banks",
		})
	} else {
		ctx.JSON(200, banks)
	}
}

func ViewBank(ctx *gin.Context, db *pg.DB) {
	nbid, _ := strconv.ParseUint(ctx.Param("bid"), 10, 0)
	bank, err := Models.ShowB(db, uint(nbid))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in printing all banks",
		})
	} else {
		ctx.JSON(200, bank)
	}
}

func UpdateBank(ctx *gin.Context, db *pg.DB) {
	var newBank Models.Bank
	if err := ctx.ShouldBindJSON(&newBank); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := Models.UpdateB(db, &newBank)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(200, newBank)
	}
}

func DeleteBank(ctx *gin.Context, db *pg.DB) {
	var newBank Models.Bank
	if err := ctx.ShouldBindJSON(&newBank); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := Models.DeleteB(db, &newBank)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in deleting banks",
		})
	} else {
		ctx.JSON(200, gin.H{
			"message": "deleted",
		})
	}
}
