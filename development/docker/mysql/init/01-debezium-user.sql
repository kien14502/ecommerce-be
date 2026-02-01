-- Create debezium user
CREATE USER IF NOT EXISTS 'debezium'@'%'
IDENTIFIED WITH mysql_native_password
BY 'dbz_pass';

-- Grant permissions for Debezium
GRANT SELECT, RELOAD, SHOW DATABASES,
      REPLICATION SLAVE,
      REPLICATION CLIENT
ON *.* TO 'debezium'@'%';

FLUSH PRIVILEGES;
