package config

import (
	"github.com/scukonick/geddis/db"
	"github.com/scukonick/geddis/server_api/go"
)

// Config represents all application config
type Config struct {
	Store     db.StoreConfig
	ServerAPI geddis.Config
}
