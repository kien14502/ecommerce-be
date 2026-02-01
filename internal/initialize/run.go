package initialize

import (
	"github.com/gin-gonic/gin"
)

func Run() *gin.Engine {
	LoadConfigInit()
	LoggerInit()
	MySqlInit()
	RedisInit()
	r := RouterInit()
	return r
}
