package api

import (
	"dogo/pkg/model"
	"dogo/pkg/repository"
	"dogo/pkg/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CommandCreateEndpoint(c *gin.Context) {
	var item model.Command
	if err := c.ShouldBind(&item); err != nil {
		panic(err)
	}

	item.ID = utils.UUID()
	item.Created = utils.NowJsonTime()

	if err := repository.CreateNewCommand(&item); err != nil {
		panic(err)
	}

	Success(c, item)
}

func CommandPagingEndpoint(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")
	content := c.Query("content")

	var items = make([]model.Command, 0)
	var total int64 = 0

	if err := repository.FindPageCommand(pageIndex, pageSize, &total, &items, name, content); err != nil {
		panic(err)
	}

	Success(c, gin.H{
		"total": total,
		"items": items,
	})
}

func CommandUpdateEndpoint(c *gin.Context) {
	id := c.Param("id")

	if err := repository.FindCommandById(&model.Command{}, id); err != nil {
		panic(err)
	}

	var item model.Command
	if err := c.ShouldBind(&item); err != nil {
		panic(err)
	}

	if err := repository.UpdateCommandById(&item, id); err != nil {
		panic(err)
	}

	Success(c, nil)
}

func CommandDeleteEndpoint(c *gin.Context) {
	id := c.Param("id")
	if err := repository.DeleteCommandById(id); err != nil {
		panic(err)
	}

	Success(c, nil)
}

func CommandGetEndpoint(c *gin.Context) {
	id := c.Param("id")

	var item model.Command
	if err := repository.FindCommandById(&item, id); err != nil {
		panic(err)
	}

	Success(c, item)
}
