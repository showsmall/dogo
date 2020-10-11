package config

import (
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

var DB *gorm.DB

var Cache *cache.Cache

var Dogo *DogoConfig
