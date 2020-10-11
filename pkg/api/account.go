package api

import (
	"dogo/pkg/config"
	"dogo/pkg/model"
	"dogo/pkg/repository"
	utils "dogo/pkg/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type LoginAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginEndpoint(c *gin.Context) {
	var loginAccount LoginAccount
	if err := c.ShouldBind(&loginAccount); err != nil {
		panic(err)
	}

	var user model.User
	if err := repository.FindUserByUsername(&user, loginAccount.Username); err != nil {
		panic("您输入的用户名或密码不正确 :(")
	}
	if err := utils.Encoder.Match([]byte(user.Password), []byte(loginAccount.Password)); err != nil {
		panic("您输入的用户名或密码不正确 :)")
	}

	token := utils.UUID()

	config.Cache.Set(token, user, time.Minute*time.Duration(30))

	Success(c, token)
}

func LogoutEndpoint(c *gin.Context) {
	token := c.GetHeader("X-Auth-Token")
	config.Cache.Delete(token)
	Success(c, nil)
}

func ChangePasswordEndpoint(c *gin.Context) {

}

func InfoEndpoint(c *gin.Context) {
	token := c.GetHeader("X-Auth-Token")
	account, _ := config.Cache.Get(token)
	Success(c, account)
}
