package model

import (
	utils "dogo/pkg/utils"
	"reflect"
)

type User struct {
	ID       string         `gorm:"primary_key" json:"id"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Nickname string         `json:"nickname"`
	Type     string         `json:"type"`
	Enabled  bool           `json:"enabled"`
	Created  utils.JsonTime `json:"created"`
}

func (r *User) TableName() string {
	return "users"
}

func (r *User) IsEmpty() bool {
	return reflect.DeepEqual(r, User{})
}

type Asset struct {
	ID           string         `gorm:"primary_key " json:"id"`
	Name         string         `json:"name"`
	IP           string         `json:"ip"`
	Protocol     string         `json:"protocol"`
	Port         int            `json:"port"`
	AccountType  string         `json:"accountType"`
	Username     string         `json:"username"`
	Password     string         `json:"password"`
	CredentialId string         `json:"credentialId"`
	PrivateKey   string         `json:"privateKey"`
	Passphrase   string         `json:"passphrase"`
	Description  string         `json:"description"`
	Active       bool           `json:"active"`
	Created      utils.JsonTime `json:"created"`
}

func (r *Asset) TableName() string {
	return "assets"
}

type Command struct {
	ID      string         `gorm:"primary_key" json:"id"`
	Name    string         `json:"name"`
	Content string         `json:"content"`
	Created utils.JsonTime `json:"created"`
}

func (r *Command) TableName() string {
	return "commands"
}

type Credential struct {
	ID       string         `gorm:"primary_key" json:"id"`
	Name     string         `json:"name"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Created  utils.JsonTime `json:"created"`
}

func (r *Credential) TableName() string {
	return "credentials"
}

type Session struct {
	ID               string         `gorm:"primary_key" json:"id"`
	Protocol         string         `json:"protocol"`
	ConnectionId     string         `json:"connectionId"`
	AssetId          string         `json:"assetId"`
	Creator          string         `json:"creator"`
	Width            int            `json:"width"`
	Height           int            `json:"height"`
	Status           string         `json:"status"`
	ConnectedTime    utils.JsonTime `json:"connectedTime"`
	DisconnectedTime utils.JsonTime `json:"disconnectedTime"`
}

func (r *Session) TableName() string {
	return "sessions"
}
