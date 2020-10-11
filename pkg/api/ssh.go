package api

import (
	"bytes"
	"dogo/pkg/model"
	"dogo/pkg/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type dogoWriter struct {
	b  bytes.Buffer
	mu sync.Mutex
}

func (w *dogoWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Write(p)
}

func (w *dogoWriter) Read() ([]byte, int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	p := w.b.Bytes()
	buf := make([]byte, len(p))
	read, err := w.b.Read(buf)
	w.b.Reset()
	if err != nil {
		return nil, 0, err
	}
	return buf, read, err
}

func SSHEndpoint(c *gin.Context) {
	ws, err := UpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("WebSocket upgrade error", err)
	}

	assetId := c.Query("assetId")
	width, _ := strconv.Atoi(c.DefaultQuery("width", "1024"))
	height, _ := strconv.Atoi(c.DefaultQuery("height", "768"))

	var asset model.Asset
	if err := repository.FindAssetById(&asset, assetId); err != nil {
		WriteMessage(ws, "获取资产失败")
		return
	}

	if asset.AccountType == "credential" {
		var credential model.Credential
		if err := repository.FindCredentialById(&credential, asset.CredentialId); err != nil {
			WriteMessage(ws, "获取资产凭证失败")
			return
		}
		asset.Username = credential.Username
		asset.Password = credential.Password
	}

	config := &ssh.ClientConfig{
		Timeout: 1 * time.Second,
		User:    asset.Username,
		Auth:    []ssh.AuthMethod{ssh.Password(asset.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := fmt.Sprintf("%s:%d", asset.IP, asset.Port)

	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		WriteMessage(ws, "建立连接失败 "+err.Error())
		return
	}

	session, err := sshClient.NewSession()
	if err != nil {
		WriteMessage(ws, "创建会话失败 "+err.Error())
		return
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", height, width, modes); err != nil {
		WriteMessage(ws, "获取pty异常 "+err.Error())
		return
	}

	var b dogoWriter
	session.Stdout = &b
	session.Stderr = &b

	stdinPipe, err := session.StdinPipe()
	if err != nil {
		WriteMessage(ws, "获取会话输入失败 "+err.Error())
		return
	}

	if err := session.Shell(); err != nil {
		WriteMessage(ws, "获取shell失败  "+err.Error())
		return
	}

	go func() {

		for true {

			p, n, err := b.Read()
			if err != nil {
				continue
			}
			if n > 0 {
				WriteByteMessage(ws, p)
			}
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
	}()

	for true {
		_, message, err := ws.ReadMessage()
		if err != nil {
			continue
		}
		_, err = stdinPipe.Write(message)
		if err != nil {
			log.Println("Tunnel write:", err)
		}
	}
}

func WriteMessage(ws *websocket.Conn, message string) {
	WriteByteMessage(ws, []byte(message))
}

func WriteByteMessage(ws *websocket.Conn, p []byte) {
	err := ws.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		log.Println("write:", err)
	}
}
