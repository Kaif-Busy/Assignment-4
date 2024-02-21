package operations

import (
	"bank/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)
type DisplayBranch struct {
	ID        uint 
	Address   string
	IFSCCode  string
	BankID    uint	

}

func AddBranch(ctx *gin.Context, db *pg.DB) {
	var newBranch Models.Branch
	if err := ctx.ShouldBindJSON(&newBranch); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := newBranch.SaveBranch(db); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Branch added successfully", "ID": newBranch.ID})
}

func ViewAllBranches(ctx *gin.Context, db *pg.DB) {
	branches, err := Models.ShowAllBranch(db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in printing all Branchs",
		})
	} else {
		var displayBranches []DisplayBranch
		for _,branch:= range branches{
			var displayBranch DisplayBranch
			displayBranch.ID=branch.ID
			displayBranch.Address=branch.Address
			displayBranch.IFSCCode=branch.IFSCCode
			displayBranch.BankID=branch.BankID
			displayBranches=append(displayBranches,displayBranch)
		}
		ctx.JSON(200, displayBranches)
	}
}

func ViewBranch(ctx *gin.Context, db *pg.DB) {
	nbid, _ := strconv.ParseUint(ctx.Param("bid"), 10, 0)
	branch, err := Models.ShowBranch(db, uint(nbid))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in printing all Branches",
		})
	} else {
			var displayBranch DisplayBranch
			displayBranch.ID=branch.ID
			displayBranch.Address=branch.Address
			displayBranch.IFSCCode=branch.IFSCCode
			displayBranch.BankID=branch.BankID
		ctx.JSON(200, displayBranch)
	}
}

func UpdateBranch(ctx *gin.Context, db *pg.DB) {
	var newBranch Models.Branch
	if err := ctx.ShouldBindJSON(&newBranch); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := Models.UpdateBranch(db, &newBranch)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(200,gin.H{"Message": "Updated"})
	}
}

func DeleteBranch(ctx *gin.Context, db *pg.DB) {
	var newBranch Models.Branch
	if err := ctx.ShouldBindJSON(&newBranch); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := Models.DeleteBranch(db, &newBranch)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in deleting Branchs",
		})
	} else {
		ctx.JSON(200, gin.H{
			"message": "deleted",
		})
	}
}
