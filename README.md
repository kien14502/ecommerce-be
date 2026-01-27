<!-- Database -->

$user: root
$password: root12345
$db_name: ecommerce_db
$docker_name :ecommerce-mysql

<!-- dump database command line -->

docker exec -it ecommerce-mysql mysqldump -uroot -proot12345 --databases ecommerce_db --add-drop-database --add-drop-table --add-drop-trigger --add-locks --no-data > migrations/ecommerce.sql
