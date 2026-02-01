<!-- Database -->

$user: root
$password: root12345
$db_name: ecommerce_db
$docker_name :ecommerce-mysql

<!-- dump database command line -->

docker exec -it ecommerce-mysql mysqldump -uroot -proot12345 --databases ecommerce_db --add-drop-database --add-drop-table --add-drop-trigger --add-locks --no-data > migrations/ecommerce.sql

<!-- kafka -->
<!-- redis -->
<!-- mailer -->
<!-- debezium -->

curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d scripts/register-mysql.json

docker run -it --rm --name connect -p 8083:8083 -e GROUP_ID=1 -e CONFIG_STORAGE_TOPIC=my_connect_configs -e OFFSET_STORAGE_TOPIC=my_connect_offsets -e STATUS_STORAGE_TOPIC=my_connect_statuses --link kafka:kafka --link ecommerce-mysql:ecommerce-mysql quay.io/debezium/connect:3.4

<!-- generate wire -->
