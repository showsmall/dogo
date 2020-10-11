package api

import (
	"dogo/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var UpGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"guacamole"},
}

func SetupRoutes() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(Cors)
	r.Use(Recover)

	r.POST("/login", LoginEndpoint)
	r.GET("/tunnel", TunnelEndpoint)
	r.GET("/ssh", SSHEndpoint)

	r.Use(Auth)
	r.POST("/logout", LogoutEndpoint)
	r.POST("/change-password", ChangePasswordEndpoint)
	r.GET("/info", InfoEndpoint)

	rUser := r.Group("/users")
	rUser.Use(Auth)
	{
		rUser.GET("", UserPagingEndpoint)
		rUser.POST("", UserCreateEndpoint)
		rUser.PUT("/:id", UserUpdateEndpoint)
		rUser.DELETE("/:id", UserDeleteEndpoint)
		rUser.GET("/:id", UserGetEndpoint)
	}

	r.GET("/assets-all", AssetAllEndpoint)
	rAsset := r.Group("/assets")
	//rUser.Use(Auth)
	{
		rAsset.GET("", AssetPagingEndpoint)
		rAsset.POST("", AssetCreateEndpoint)
		rAsset.PUT("/:id", AssetUpdateEndpoint)
		rAsset.DELETE("/:id", AssetDeleteEndpoint)
		rAsset.GET("/:id", AssetGetEndpoint)
	}

	rCommand := r.Group("/commands")
	//rUser.Use(Auth)
	{
		rCommand.GET("", CommandPagingEndpoint)
		rCommand.POST("", CommandCreateEndpoint)
		rCommand.PUT("/:id", CommandUpdateEndpoint)
		rCommand.DELETE("/:id", CommandDeleteEndpoint)
		rCommand.GET("/:id", CommandGetEndpoint)
	}

	rCredential := r.Group("/credentials")
	r.GET("/credentials-all", CredentialAllEndpoint)
	{
		rCredential.GET("", CredentialPagingEndpoint)
		rCredential.POST("", CredentialCreateEndpoint)
		rCredential.PUT("/:id", CredentialUpdateEndpoint)
		rCredential.DELETE("/:id", CredentialDeleteEndpoint)
		rCredential.GET("/:id", CredentialGetEndpoint)
	}

	rSession := r.Group("/sessions")
	//rUser.Use(Auth)
	{
		rSession.GET("", SessionPagingEndpoint)
		rSession.POST("", SessionCreateEndpoint)
		//rSession.POST("/:id/discontent", SessionDiscontentEndpoint)
		rSession.DELETE("/:id", SessionDeleteEndpoint)
	}

	return r
}

func Fail(c *gin.Context, code int, message string) {
	c.JSON(200, gin.H{
		"code":    code,
		"message": message,
	})
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code":    1,
		"message": "success",
		"data":    data,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(200, gin.H{
		"code":    -1,
		"message": message,
	})
}

func GetCurrentAccount(c *gin.Context) (interface{}, bool) {
	token := c.GetHeader("X-Auth-Token")
	return config.Cache.Get(token)
}
