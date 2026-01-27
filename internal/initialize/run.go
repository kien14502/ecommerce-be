package initialize

import (
	"strconv"

	"github.com/kien14502/ecommerce-be/global"
)

func Run() {
	LoadConfigInit()
	LoggerInit()
	MySqlInit()
	RedisInit()
	r := RouterInit()

	port := global.Config.Server.Port
	r.Run(":" + strconv.Itoa(port))
}
