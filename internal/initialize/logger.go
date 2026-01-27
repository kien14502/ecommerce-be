package initialize

import (
	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/pkg/logger"
)

func LoggerInit() {
	global.Logger = logger.NewLogger(global.Config.Logger)
}
