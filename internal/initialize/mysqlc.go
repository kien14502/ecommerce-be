package initialize

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kien14502/ecommerce-be/global"
	"go.uber.org/zap"
)

func checkErrPanic(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString, zap.Error(err))
		panic(err)
	}
}

func MySqlInitc() {
	config := global.Config.MySql
	dsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	var configDSN = fmt.Sprintf(dsn, config.Username, config.Password, config.Host, config.Port, config.DBName)
	db, err := sql.Open("mysql", configDSN)
	checkErrPanic(err, "Failed to connect to database")
	global.Logger.Info("Connect my sql successfully")
	global.Mdbc = db
	db.Exec(fmt.Sprintf("USE `%s`", config.DBName))
}
