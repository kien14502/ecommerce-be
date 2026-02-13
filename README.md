<!-- Database -->

$user: root
$password: root12345
$db_name: ecommerce_db
$docker_name :ecommerce-mysql

<!-- gorm -->

// func SetPoolc() {
// mySql, err := global
// checkErrPanic(err, "Failed to get database instance")
// config := global.Config.MySql
// mySql.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime))
// mySql.SetMaxOpenConns(config.MaxOpen)
// mySql.SetConnMaxIdleTime(time.Duration(config.MaxIdle))
// }

// func migrateTablesc() {
// err := global.Mdb.AutoMigrate(
// &models.User{},
// &models.RoleModel{},
// )

// checkErrPanic(err, "Failed to migrate database")
// }

// func genTableDaoc() {
// g := gen.NewGenerator(gen.Config{
// OutPath: "internal/models",
// Mode: gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
// })

// g.UseDB(global.Mdb)
// g.Execute()
// }

<!-- dump database command line -->

docker exec -it ecommerce-mysql mysqldump -uroot -proot12345 --databases ecommerce_db --add-drop-database --add-drop-table --add-drop-trigger --add-locks --no-data > migrations/ecommerce.sql

<!-- kafka -->
<!-- redis -->
<!-- mailer -->
<!-- debezium -->

curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d scripts/register-mysql.json

docker run -it --rm --name connect -p 8083:8083 -e GROUP_ID=1 -e CONFIG_STORAGE_TOPIC=my_connect_configs -e OFFSET_STORAGE_TOPIC=my_connect_offsets -e STATUS_STORAGE_TOPIC=my_connect_statuses --link kafka:kafka --link ecommerce-mysql:ecommerce-mysql quay.io/debezium/connect:3.4

<!-- generate wire -->
<!-- goose -->

# Command:

    - Create table: goose -dir sql/schema create ${TABLE_NAME} sql
    - Generate: sqlc generate
