package api

import (
	"dogo/pkg/config"
	"dogo/pkg/guacd"
	"dogo/pkg/model"
	"dogo/pkg/repository"
	"dogo/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

func TunnelEndpoint(c *gin.Context) {

	ws, err := UpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("WebSocket upgrade error", err)
	}

	width := c.DefaultQuery("width", "1024")
	height := c.DefaultQuery("height", "768")
	sessionId := c.Query("sessionId")
	connectionId := c.Query("connectionId")

	configuration := guacd.NewConfiguration()
	configuration.SetParameter("width", width)
	configuration.SetParameter("height", height)

	var session model.Session
	if len(connectionId) > 0 {
		if err := repository.FindSessionByConnectionId(&session, connectionId); err != nil {
			return
		}
		configuration.ConnectionID = connectionId
	} else {
		if err := repository.FindSessionById(&session, sessionId); err != nil {
			return
		}
		var asset model.Asset
		if err := repository.FindAssetById(&asset, session.AssetId); err != nil {
			return
		}

		if asset.AccountType == "credential" {
			var credential model.Credential
			if err := repository.FindCredentialById(&credential, asset.CredentialId); err != nil {
				log.Fatal("获取授权凭证失败")
				return
			}
			asset.Username = credential.Username
			asset.Password = credential.Password
		}

		configuration.Protocol = asset.Protocol
		switch asset.Protocol {
		case "rdp":
			configuration.SetParameter("username", asset.Username)
			configuration.SetParameter("password", asset.Password)

			configuration.SetParameter("security", "any")
			configuration.SetParameter("ignore-cert", "true")
			configuration.SetParameter("enable-drive", "true")
			configuration.SetParameter("drive-name", config.Dogo.Guacd.RDP.DriveName)
			configuration.SetParameter("drive-path", config.Dogo.Guacd.RDP.DrivePath)
			configuration.SetParameter("create-drive-path", "true")

		case "ssh":
			configuration.SetParameter("username", asset.Username)
			configuration.SetParameter("password", asset.Password)

			configuration.SetParameter("enable-sftp", "true")
			configuration.SetParameter("font-name", config.Dogo.Guacd.SSH.FontName)
			configuration.SetParameter("font-size", config.Dogo.Guacd.SSH.FontSize)
			configuration.SetParameter("color-scheme", config.Dogo.Guacd.SSH.ColorScheme)
		case "vnc":
			configuration.SetParameter("autoretry", config.Dogo.Guacd.VNC.Autoretry)
			configuration.SetParameter("password", asset.Password)
		case "telnet":
			configuration.SetParameter("username", asset.Username)
			configuration.SetParameter("password", asset.Password)
		}

		configuration.SetParameter("hostname", asset.IP)
		configuration.SetParameter("port", strconv.Itoa(asset.Port))

		//configuration.SetParameter("recording-path", "")
		//configuration.SetParameter("create-recording-path", "true")
		//configuration.SetParameter("recording-name", sessionId+".guac")
	}

	tunnel := guacd.NewTunnel(config.Dogo.Guacd.Addr, configuration)

	if len(session.ConnectionId) == 0 {
		session.ConnectionId = tunnel.UUID
		session.ConnectedTime = utils.NowJsonTime()
		if err := repository.UpdateSessionById(&session, sessionId); err != nil {
			log.Fatal("更新会话失败")
			return
		}
	}
	go func() {
		for true {
			instruction, err := tunnel.Read()
			if err != nil {
				log.Println("Tunnel read:", err)
				continue
			}
			err = ws.WriteMessage(websocket.TextMessage, instruction)
			if err != nil {
				log.Println("write:", err)
			}
		}
	}()

	for true {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			continue
		}
		_, err = tunnel.Write(message)
		if err != nil {
			log.Println("Tunnel write:", err)
		}
	}
}
