package initialize

import (
	"github.com/gin-gonic/gin"
)

func Run() *gin.Engine {
	LoadConfigInit()
	LoggerInit()
	MySqlInitc()
	RedisInit()
	InitKafka()
	r := RouterInit()
	return r
}
