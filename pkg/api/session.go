package api

import (
	"dogo/pkg/config"
	"dogo/pkg/model"
	"dogo/pkg/repository"
	utils "dogo/pkg/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func SessionPagingEndpoint(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")

	var items = make([]model.Session, 0)
	var total int64 = 0

	if err := repository.FindPageSession(pageIndex, pageSize, &total, &items, status); err != nil {
		panic(err)
	}

	Success(c, gin.H{
		"total": total,
		"items": items,
	})
}

func SessionDeleteEndpoint(c *gin.Context) {
	id := c.Param("id")
	if err := repository.DeleteSessionById(id); err != nil {
		panic(err)
	}

	Success(c, nil)
}

func SessionDiscontentEndpoint(c *gin.Context) {

}

func SessionCreateEndpoint(c *gin.Context) {
	assetId := c.Query("assetId")

	var item model.Asset
	if err := repository.FindAssetById(&item, assetId); err != nil {
		panic(err)
	}

	session := &model.Session{
		ID:       utils.UUID(),
		AssetId:  item.ID,
		Protocol: item.Protocol,
		Status:   config.NotContented,
	}

	if err := repository.CreateNewSession(session); err != nil {
		panic(err)
	}

	Success(c, session)
}
