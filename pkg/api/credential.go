package api

import (
	"dogo/pkg/model"
	"dogo/pkg/repository"
	"dogo/pkg/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CredentialAllEndpoint(c *gin.Context) {
	var items = make([]model.Credential, 0)
	if err := repository.FindAllCredential(&items); err != nil {
		panic(err)
	}
	Success(c, items)
}
func CredentialCreateEndpoint(c *gin.Context) {
	var item model.Credential
	if err := c.ShouldBind(&item); err != nil {
		panic(err)
	}

	item.ID = utils.UUID()
	item.Created = utils.NowJsonTime()

	if err := repository.CreateNewCredential(&item); err != nil {
		panic(err)
	}

	Success(c, item)
}

func CredentialPagingEndpoint(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")

	var items = make([]model.Credential, 0)
	var total int64 = 0

	if err := repository.FindPageCredential(pageIndex, pageSize, &total, &items, name); err != nil {
		panic(err)
	}

	Success(c, gin.H{
		"total": total,
		"items": items,
	})
}

func CredentialUpdateEndpoint(c *gin.Context) {
	id := c.Param("id")

	if err := repository.FindCredentialById(&model.Credential{}, id); err != nil {
		panic(err)
	}

	var item model.Credential
	if err := c.ShouldBind(&item); err != nil {
		panic(err)
	}

	if err := repository.UpdateCredentialById(&item, id); err != nil {
		panic(err)
	}

	Success(c, nil)
}

func CredentialDeleteEndpoint(c *gin.Context) {
	id := c.Param("id")
	if err := repository.DeleteCredentialById(id); err != nil {
		panic(err)
	}

	Success(c, nil)
}

func CredentialGetEndpoint(c *gin.Context) {
	id := c.Param("id")

	var item model.Credential
	if err := repository.FindCredentialById(&item, id); err != nil {
		panic(err)
	}

	Success(c, item)
}
