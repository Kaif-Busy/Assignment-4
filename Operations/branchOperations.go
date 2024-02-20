package operations

import (
	"bank/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

func AddBranch(ctx *gin.Context, db *pg.DB) {
	var newBranch Models.Branch
	if err := ctx.ShouldBindJSON(&newBranch); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := newBranch.SaveBr(db); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Branch added successfully", "data": newBranch})
}

func ViewAllBranches(ctx *gin.Context, db *pg.DB) {
	branches, err := Models.ShowAllBr(db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in printing all Branchs",
		})
	} else {
		ctx.JSON(200, branches)
	}
}

func ViewBranch(ctx *gin.Context, db *pg.DB) {
	nbid, _ := strconv.ParseUint(ctx.Param("bid"), 10, 0)
	branch, err := Models.ShowBr(db, uint(nbid))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in printing all Branches",
		})
	} else {
		ctx.JSON(200, branch)
	}
}

func UpdateBranch(ctx *gin.Context, db *pg.DB) {
	var newBranch Models.Branch
	if err := ctx.ShouldBindJSON(&newBranch); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := Models.UpdateBr(db, &newBranch)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(200, newBranch)
	}
}

func DeleteBranch(ctx *gin.Context, db *pg.DB) {
	var newBranch Models.Branch
	if err := ctx.ShouldBindJSON(&newBranch); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := Models.DeleteBr(db, &newBranch)

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
