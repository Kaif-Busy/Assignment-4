package operations

import (
	"bank/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)
type DisplayAllBanks struct{
	ID uint
	Name string
}

type DisplayBank struct{
	ID uint
	Name string
	Branches []string
}

func AddBank(ctx *gin.Context, db *pg.DB) {
	var newBank Models.Bank
	if err := ctx.ShouldBindJSON(&newBank); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := newBank.SaveBank(db); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to save Bank to database"})
		return
	}
	ctx.JSON(200, gin.H{"message": "Bank added successfully", "ID": newBank.ID})
}

func ViewAllBanks(ctx *gin.Context, db *pg.DB) {
	banks, err := Models.ShowAllBank(db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in printing all banks",
		})
	} else {
		var allBanks []DisplayAllBanks
		for _,bank:=range banks{
			var allBank DisplayAllBanks
			allBank.ID=bank.ID
			allBank.Name=bank.BankName
			allBanks=append(allBanks, allBank)
		}
		
		ctx.JSON(200, allBanks)
	}
}

func ViewBank(ctx *gin.Context, db *pg.DB) {
	nbid, _ := strconv.ParseUint(ctx.Param("bid"), 10, 0)
	bank, err := Models.ShowBank(db, uint(nbid))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in printing all banks",
		})
	} else {
		var displayBank DisplayBank
		displayBank.Name=bank.BankName
		displayBank.ID=bank.ID
		for _,branch:= range bank.Branches{
			displayBank.Branches=append(displayBank.Branches, branch.Address)
		}

		ctx.JSON(200, displayBank)
	}
}

func UpdateBank(ctx *gin.Context, db *pg.DB) {
	var newBank Models.Bank
	if err := ctx.ShouldBindJSON(&newBank); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := Models.UpdateBank(db, &newBank)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(200,gin.H{"message":"Updated"})
	}
}

func DeleteBank(ctx *gin.Context, db *pg.DB) {
	var newBank Models.Bank
	if err := ctx.ShouldBindJSON(&newBank); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := Models.DeleteBank(db, &newBank)

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
