package server

import (
	"github.com/gin-gonic/gin"

	"github.com/figment-networks/mina-indexer/config"
)

// SetGinDefaults changes Gin behavior base on application environment
func SetGinDefaults(cfg *config.Config) {
	if cfg.IsProduction() {
		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)
	}
}
