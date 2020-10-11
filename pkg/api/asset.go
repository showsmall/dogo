package api

import (
	"dogo/pkg/model"
	"dogo/pkg/repository"
	"dogo/pkg/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AssetCreateEndpoint(c *gin.Context) {
	var item model.Asset
	if err := c.ShouldBind(&item); err != nil {
		panic(err)
	}

	item.ID = utils.UUID()
	item.Created = utils.NowJsonTime()

	if err := repository.CreateNewAsset(&item); err != nil {
		panic(err)
	}

	Success(c, item)
}

func AssetPagingEndpoint(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")

	var items = make([]model.Asset, 0)
	var total int64 = 0

	if err := repository.FindPageAsset(pageIndex, pageSize, &total, &items, name); err != nil {
		panic(err)
	}

	Success(c, gin.H{
		"total": total,
		"items": items,
	})
}

func AssetAllEndpoint(c *gin.Context) {
	var items = make([]model.Asset, 0)
	if err := repository.FindAllAsset(&items); err != nil {
		panic(err)
	}
	Success(c, items)
}

func AssetUpdateEndpoint(c *gin.Context) {
	id := c.Param("id")

	if err := repository.FindAssetById(&model.Asset{}, id); err != nil {
		panic(err)
	}

	var item model.Asset
	if err := c.ShouldBind(&item); err != nil {
		panic(err)
	}

	if err := repository.UpdateAssetById(&item, id); err != nil {
		panic(err)
	}

	Success(c, nil)
}

func AssetDeleteEndpoint(c *gin.Context) {
	id := c.Param("id")
	if err := repository.DeleteAssetById(id); err != nil {
		panic(err)
	}

	Success(c, nil)
}

func AssetGetEndpoint(c *gin.Context) {
	id := c.Param("id")

	var item model.Asset
	if err := repository.FindAssetById(&item, id); err != nil {
		panic(err)
	}

	Success(c, item)
}
