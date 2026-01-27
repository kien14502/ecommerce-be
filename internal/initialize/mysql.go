package initialize

import (
	"fmt"
	"time"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/models"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func checkErrPanic(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString, zap.Error(err))
		panic(err)
	}
}

func MySqlInit() {
	config := global.Config.MySql
	dsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	var configDSN = fmt.Sprintf(dsn, config.Username, config.Password, config.Host, config.Port, config.DBName)
	db, err := gorm.Open(mysql.Open(configDSN), &gorm.Config{
		SkipDefaultTransaction: false,
	})
	checkErrPanic(err, "Failed to connect to database")
	global.Logger.Info("Connect my sql successfully")
	global.Mdb = db
	db.Exec(fmt.Sprintf("USE `%s`", config.DBName))

	// set connection pool
	SetPool()
	migrateTables()
	genTableDao()
}

func SetPool() {
	mySql, err := global.Mdb.DB()
	checkErrPanic(err, "Failed to get database instance")
	config := global.Config.MySql
	mySql.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime))
	mySql.SetMaxOpenConns(config.MaxOpen)
	mySql.SetConnMaxIdleTime(time.Duration(config.MaxIdle))
}

func migrateTables() {
	err := global.Mdb.AutoMigrate(
		&models.UserModel{},
		&models.Role{},
	)
	checkErrPanic(err, "Failed to migrate database")
}

func genTableDao() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "internal/models",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(global.Mdb)
	g.Execute()
}
