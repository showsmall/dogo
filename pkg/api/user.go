package api

import (
	"dogo/pkg/model"
	"dogo/pkg/repository"
	"dogo/pkg/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func UserCreateEndpoint(c *gin.Context) {
	var item model.User
	if err := c.ShouldBind(&item); err != nil {
		panic(err)
	}

	var pass []byte
	var err error
	if pass, err = utils.Encoder.Encode([]byte("admin")); err != nil {
		return
	}
	item.Password = string(pass)

	item.ID = utils.UUID()
	item.Created = utils.NowJsonTime()

	if err := repository.CreateNewUser(&item); err != nil {
		panic(err)
	}

	Success(c, item)
}

func UserPagingEndpoint(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var items = make([]model.User, 0)
	var total int64 = 0

	if err := repository.FindPageUser(pageIndex, pageSize, &total, &items); err != nil {
		panic(err)
	}

	Success(c, gin.H{
		"total": total,
		"items": items,
	})
}

func UserUpdateEndpoint(c *gin.Context) {
	id := c.Param("id")

	if err := repository.FindUserById(&model.User{}, id); err != nil {
		panic(err)
	}

	var item model.User
	if err := c.ShouldBind(&item); err != nil {
		panic(err)
	}

	if err := repository.UpdateUserById(&item, id); err != nil {
		panic(err)
	}

	Success(c, nil)
}

func UserDeleteEndpoint(c *gin.Context) {
	id := c.Param("id")
	if err := repository.DeleteUserById(id); err != nil {
		panic(err)
	}

	Success(c, nil)
}

func UserGetEndpoint(c *gin.Context) {
	id := c.Param("id")

	var item model.User
	if err := repository.FindUserById(&item, id); err != nil {
		panic(err)
	}

	Success(c, item)
}
